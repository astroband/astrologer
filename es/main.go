package es

import (
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
)

var ledgerHeaderIndexName = "ledger"
var txIndexName = "tx"
var opIndexName = "op"
var balanceIndexName = "balance"
var tradesIndexName = "trades"
var signerHistoryIndexName = "signers"

// Indexable represents object that can be indexed for ElasticSearch
type Indexable interface {
	DocID() *string
	IndexName() string
}

func fatalIfError(res *esapi.Response, err error) {
	if err != nil {
		log.Fatal(err)
	}

	if res.IsError() {
		log.Fatal(res)
	}
}
