package es

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/stellar/go/xdr"
)

var ledgerHeaderIndexName = "ledger"
var txIndexName = "tx"
var opIndexName = "op"
var balanceIndexName = "balance"

// Asset represents es-serializable asset
type Asset struct {
	Code   string `json:"code"`
	Issuer string `json:"issuer,omitempty"`
	Key    string `json:"key"`
}

func NewNativeAsset() *Asset {
	return &Asset{"native", "", "native"}
}

func NewAsset(a *xdr.Asset) *Asset {
	var t, c, i string

	a.MustExtract(&t, &c, &i)

	if t == "native" {
		return NewNativeAsset()
	}

	return &Asset{c, i, fmt.Sprintf("%s-%s", c, i)}
}

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
