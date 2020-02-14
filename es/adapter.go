package es

import (
	"bytes"
	"encoding/json"
	goES "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Adapter interface {
	MinMaxSeq() (min, max int)
	LedgerSeqRangeQuery(ranges []map[string]interface{}) map[string]interface{}
	GetLedgerSeqsInRange(min, max int) []int
	LedgerCountInRange(min, max int) int
	IndexExists(name string) bool
	CreateIndex(name, body string)
	DeleteIndex(name string)
	BulkInsert(payload *bytes.Buffer) (success bool)
	IndexWithRetries(payload *bytes.Buffer, retriesCount int)
}

type Client struct {
	rawClient *goES.Client
}

func (es *Client) IndexExists(name string) bool {
	res, err := es.rawClient.Indices.Get([]string{name})

	if err != nil {
		log.Fatal(err)
	}

	return res.StatusCode != http.StatusNotFound
}

func (es *Client) DeleteIndex(name string) {
	res, err := es.rawClient.Indices.Delete([]string{name})
	fatalIfError(res, err)
}

func (es *Client) CreateIndex(name, body string) {
	create := es.rawClient.Indices.Create

	res, err := create(
		name,
		create.WithBody(strings.NewReader(body)),
		create.WithIncludeTypeName(false),
	)
	fatalIfError(res, err)
}

func (es *Client) searchLedgers(query map[string]interface{}) (r map[string]interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.rawClient.Search(
		es.rawClient.Search.WithIndex("ledger"),
		es.rawClient.Search.WithBody(&buf),
	)

	fatalIfError(res, err)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	res.Body.Close()

	return r
}

func (es *Client) MinMaxSeq() (min, max int) {
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"seq_stats": map[string]interface{}{
				"stats": map[string]interface{}{
					"field": "seq",
				},
			},
		},
	}

	r := es.searchLedgers(query)

	aggs := r["aggregations"].(map[string]interface{})["seq_stats"].(map[string]interface{})

	min = int(aggs["min"].(float64))
	max = int(aggs["max"].(float64))

	return min, max
}

func (es *Client) LedgerSeqRangeQuery(ranges []map[string]interface{}) map[string]interface{} {
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"seq_ranges": map[string]interface{}{
				"range": map[string]interface{}{
					"field":  "seq",
					"ranges": ranges,
				},
			},
		},
	}

	r := es.searchLedgers(query)
	aggs := r["aggregations"].(map[string]interface{})["seq_ranges"].(map[string]interface{})

	return aggs
}

func (es *Client) BulkInsert(payload *bytes.Buffer) (success bool) {
	res, err := es.rawClient.Bulk(bytes.NewReader(payload.Bytes()))

	if res != nil {
		defer res.Body.Close()
	}

	return err == nil && (res == nil || !res.IsError())
}

func (es *Client) LedgerCountInRange(min, max int) int {
	var r map[string]interface{}
	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"seq": map[string]interface{}{
					"gte": min,
					"lte": max,
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.rawClient.Count(
		es.rawClient.Count.WithIndex("ledger"),
		es.rawClient.Count.WithBody(&buf),
	)

	fatalIfError(res, err)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	res.Body.Close()

	return int(r["count"].(float64))
}

func (es *Client) GetLedgerSeqsInRange(min, max int) (seqs []int) {
	query := map[string]interface{}{
		"_source": []string{"seq"},
		"size":    max - min + 1,
		"sort": []map[string]interface{}{{
			"seq": "asc",
		}},
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"seq": map[string]interface{}{
					"gte": min,
					"lt":  max,
				},
			},
		},
	}

	r := es.searchLedgers(query)

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})
		seqs = append(seqs, int(source["seq"].(float64)))
	}

	return
}

func (es *Client) IndexWithRetries(payload *bytes.Buffer, retryCount int) {
	isIndexed := es.BulkInsert(payload)

	if !isIndexed {
		if retryCount-1 == 0 {
			log.Fatal("Retries for bulk failed, aborting")
		}

		delay := time.Duration((rand.Intn(10) + 5))
		time.Sleep(delay * time.Second)

		es.IndexWithRetries(payload, retryCount-1)
	}
}

func Connect(url string) *Client {
	esCfg := goES.Config{
		Addresses: []string{url},
	}

	client, err := goES.NewClient(esCfg)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{rawClient: client}
}

func fatalIfError(res *esapi.Response, err error) {
	if err != nil {
		log.Fatal(err)
	}

	if res.IsError() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		log.Fatal("Error in response", buf.String())
	}
}
