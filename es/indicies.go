package es

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gzigzigzeo/stellar-core-export/config"
)

const ledgerIndex = `
	{
		"mappings": {
			"_doc": {
				"properties": {
					"hash": { "type": "keyword", "index": true },
					"prev_hash": { "type": "keyword", "index": false },
					"bucket_list_hash": { "type": "keyword", "index": false },
					"seq": { "type": "long" },
					"close_time": { "type": "date" },
					"version": { "type": "long" },
					"total_coins": { "type": "long" },
					"fee_pool": { "type": "long" },
					"id_pool": { "type": "long" },
					"base_fee": { "type": "long" },
					"base_reserve": { "type": "long" },
					"max_tx_set_size": { "type": "long" }					
				}
			}
		}
	}
`

const txIndex = `
	{
		"mappings": {
			"_doc": {
				"properties": {
					"id": { "type": "keyword", "index": true },
					"idx": { "type": "integer" },
					"seq": { "type": "long" },
					"order": { "type": "keyword", "index": true },
					"fee": { "type": "long" },
					"fee_paid": { "type": "long" },
					"operation_count": { "type": "byte" },
					"close_time": { "type": "date" },
					"successful": { "type": "boolean" },
					"result_code": { "type": "byte" },
					"source_account_id": { "type": "keyword", "index": true }
				}
			}
		}
	}
`

// CreateIndicies creates all indicies in ElasticSearch database
func CreateIndicies() {
	if *config.ForceRecreateIndexes {
		deleteIndex([]string{"ledger", "tx"})
	}

	createIndex("ledger", ledgerIndex)
	createIndex("tx", txIndex)
}

func deleteIndex(index []string) {
	req := esapi.IndicesDeleteRequest{
		Index: index,
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	checkErr(res, err)
}

func createIndex(index string, body string) {
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(body),
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	checkErr(res, err)
}
