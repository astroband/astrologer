package es

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
