package es

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gzigzigzeo/stellar-core-export/config"
)

const ledgerHeaderIndex = `
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
					"fee_charged": { "type": "long" },
					"operation_count": { "type": "byte" },
					"close_time": { "type": "date" },
					"successful": { "type": "boolean" },
					"result_code": { "type": "byte" },
					"source_account_id": { "type": "keyword", "index": true },
					"memo": {
						"properties": {
							"type": { "type": "byte" },
							"value": { "type": "keyword" }
						}
					}
				}
			}
		}
	}
`

// TODO: PathPayment path, SetOptions options, DataEntry name/value
const opIndex = `
	{
		"mappings": {
			"_doc": {
				"properties": {
					"tx_id": { "type": "keyword", "index": true },
					"tx_idx": { "type": "integer" },
					"idx": { "type": "integer" },
					"seq": { "type": "long" },
					"order": { "type": "keyword", "index": true },
					"close_time": { "type": "date" },
					"successful": { "type": "boolean" },
					"result_code": { "type": "byte" },
					"tx_source_account_id": { "type": "keyword", "index": true },
					"memo": {
						"properties": {
							"type": { "type": "byte" },
							"value": { "type": "keyword" }
						}
					},
					"type": { "type": "keyword", "index": true },
					"source_account_id": { "type": "keyword", "index": true },
					"source_asset": { "type": "keyword" },
					"source_amount": { "type": "long" },
					"destination_account_id": { "type": "keyword", "index": true },
					"destination_asset": { "type": "keyword" },
					"destination_amount": { "type": "long" },
					"starting_balance": { "type": "long" },
					"offer_price": { "type": "long" },
					"offer_id": { "type": "long" },
					"trust_limit": { "type": "long" },
					"authorize": { "type": "boolean" },
					"bump_to": { "type": "long" }
				}
			}
		}
	}
`

// CreateIndicies creates all indicies in ElasticSearch database
func CreateIndicies() {
	refreshIndex(ledgerHeaderIndexName, ledgerHeaderIndex)
	refreshIndex(txIndexName, txIndex)
	refreshIndex(opIndexName, opIndex)
}

func refreshIndex(name string, body string) {
	req := esapi.IndicesGetRequest{
		Index: []string{name},
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusNotFound {
		createIndex(name, body)
		log.Printf("%s index created!", name)
	} else {
		if *config.ForceRecreateIndexes {
			deleteIndex(name)
			createIndex(name, body)
			log.Printf("%s index recreated!", name)
		} else {
			log.Printf("%s index found, skipping...", name)
		}
	}
}

func deleteIndex(index string) {
	req := esapi.IndicesDeleteRequest{
		Index: []string{index},
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	fatalIfError(res, err)
}

func createIndex(index string, body string) {
	req := esapi.IndicesCreateRequest{
		Index: index,
		Body:  strings.NewReader(body),
	}

	res, err := req.Do(context.Background(), config.ES)
	defer res.Body.Close()

	fatalIfError(res, err)
}
