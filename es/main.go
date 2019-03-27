package es

import (
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
)

var ledgerIndexName = "ledger"
var txIndexName = "tx"

func checkErr(res *esapi.Response, err error) {
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	if res.IsError() {
		log.Println(res)
		log.Fatalf("[%s] Error occured!", res.Status())
	}
}
