package es

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// IndexExists checks if an index with a given name exists in the ES cluster
func (es *Client) IndexExists(name IndexName) bool {
	res, err := es.rawClient.Indices.Get([]string{string(name)})

	if err != nil {
		log.Fatal(err)
	}

	return res.StatusCode != http.StatusNotFound
}

// DeleteIndex deletes the index from the ES cluster
func (es *Client) DeleteIndex(name IndexName) {
	res, err := es.rawClient.Indices.Delete([]string{string(name)})
	fatalIfError(res, err)
}

// CreateIndex creates an index with the given name and definition in the ES cluster
func (es *Client) CreateIndex(name IndexName, body IndexDefinition) {
	create := es.rawClient.Indices.Create

	res, err := es.rawClient.Indices.Create(
		string(name),
		create.WithBody(strings.NewReader(string(body))),
		create.WithIncludeTypeName(false),
	)
	fatalIfError(res, err)
}

func (es *Client) search(query map[string]interface{}, index IndexName) (r map[string]interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.rawClient.Search(
		es.rawClient.Search.WithIndex(string(index)),
		es.rawClient.Search.WithBody(&buf),
	)

	fatalIfError(res, err)

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	res.Body.Close()

	return r
}

func (es *Client) MinMaxSeq() (min, max int, storageEmpty bool) {
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"seq_stats": map[string]interface{}{
				"stats": map[string]interface{}{
					"field": "seq",
				},
			},
		},
	}

	r := es.search(query, ledgerHeaderIndexName)

	aggs := r["aggregations"].(map[string]interface{})["seq_stats"].(map[string]interface{})

	if aggs["min"] != nil && aggs["max"] != nil {
		min = int(aggs["min"].(float64))
		max = int(aggs["max"].(float64))
		storageEmpty = false
	} else {
		min = 0
		max = 0
		storageEmpty = true
	}

	return
}

// LedgerSeqRangeQuery fetches ledger ranges from ES
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

	r := es.search(query, ledgerHeaderIndexName)
	aggs := r["aggregations"].(map[string]interface{})["seq_ranges"].(map[string]interface{})

	return aggs
}

// BulkInsert sends the payload to ES using bulk operation
func (es *Client) BulkInsert(payload *bytes.Buffer) (success bool) {
	res, err := es.rawClient.Bulk(bytes.NewReader(payload.Bytes()))

	if res != nil {
		defer res.Body.Close()
	}

	return err == nil && (res == nil || !res.IsError())
}

// LedgerCountInRange counts number of ledgers from the given range persisted into ES
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

// GetLedgerSeqsInRange rerutns seqnums of ledgers from the given range persisted in the ES cluster
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
					"lte": max,
				},
			},
		},
	}

	r := es.search(query, ledgerHeaderIndexName)

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})
		seqs = append(seqs, int(source["seq"].(float64)))
	}

	return
}

// IndexWithRetries performs a bulk insert into ES cluster with retries on failures
func (es *Client) IndexWithRetries(payload *bytes.Buffer, retryCount int) {
	isIndexed := es.BulkInsert(payload)

	if !isIndexed {
		if retryCount-1 == 0 {
			log.Fatal("Retries for bulk failed, aborting")
		}

		log.Println("Failed, retrying...")

		delay := time.Duration((rand.Intn(10) + 5))
		time.Sleep(delay * time.Second)

		es.IndexWithRetries(payload, retryCount-1)
	}
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
