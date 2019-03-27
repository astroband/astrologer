package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gzigzigzeo/stellar-core-export/config"
)

func LedgerHeaderSerialize(h *LedgerHeader) string {
	j, err := json.Marshal(h)
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
}

func LedgerHeaderSerializeForBulk(h *LedgerHeader) string {
	return fmt.Sprintf(`{ "index": { "_index": "%s", "_id": "%s", "_type": "_doc" } }`, ledgerIndexName, h.DocID) +
		"\n" +
		LedgerHeaderSerialize(h) +
		"\n"
}

// IndexLedgerHeader Indexes ledger header
func LedgerHeaderIndex(h *LedgerHeader) {
	index(ledgerIndexName, h.DocID, LedgerHeaderSerialize(h))
}

func TransactionSerialize(t *Transaction) string {
	j, err := json.Marshal(t)
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
}

func TransactionSerializeForBulk(t *Transaction) string {
	return fmt.Sprintf(`{ "index": { "_index": "%s", "_id": "%s", "_type": "_doc" } }`, txIndexName, t.DocID) +
		"\n" +
		TransactionSerialize(t) +
		"\n"
}

// IndexLedgerHeader Indexes ledger header
func TransactionIndex(t *Transaction) {
	index(txIndexName, t.DocID, TransactionSerialize(t))
}

func BulkIndex(bulk string) {
	req := esapi.BulkRequest{
		Body:    strings.NewReader(bulk),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	checkErr(res, err)
}

func index(index string, id string, body string) {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	checkErr(res, err)
}
