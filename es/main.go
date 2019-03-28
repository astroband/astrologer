package es

import (
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
)

var ledgerHeaderIndexName = "ledger"
var txIndexName = "tx"

// Indexable represents object that can be indexed for ElasticSearch
type Indexable interface {
	DocID() string
	IndexName() string
}

func checkErr(res *esapi.Response, err error) {
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Println(res)
		log.Fatalf("[%s] Error occured!", res.Status())
	}
}
