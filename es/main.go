package es

import (
	"bytes"
	"log"

	goES "github.com/elastic/go-elasticsearch/v7"
)

// Indexable represents object that can be indexed for ElasticSearch
type Indexable interface {
	DocID() *string
	IndexName() IndexName
}

// Adapter represents the ledger storage backend
type Adapter interface {
	MinMaxSeq() (min, max int, empty bool)
	LedgerSeqRangeQuery(ranges []map[string]interface{}) map[string]interface{}
	GetLedgerSeqsInRange(min, max int) []int
	LedgerCountInRange(min, max int) int
	IndexExists(name IndexName) bool
	CreateIndex(name IndexName, body IndexDefinition)
	DeleteIndex(name IndexName)
	BulkInsert(payload *bytes.Buffer) (success bool)
	IndexWithRetries(payload *bytes.Buffer, retriesCount int)
}

// Client is a wrapper type around ElasticSearch raw client
type Client struct {
	rawClient *goES.Client
}

// Connect creates a Client configured to work with the ElasticSearch cluster
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
