package es

import (
	"log"
	"net/http"
	"strings"

	"github.com/gzigzigzeo/stellar-core-export/config"
)

const ledgerHeaderIndex = `
	{
		"mappings": {
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
`

const txIndex = `
	{
		"mappings": {
			"properties": {
				"id": { "type": "keyword", "index": true },
				"idx": { "type": "integer" },
				"seq": { "type": "long" },
				"order": { "type": "long", "index": true },
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
`

const opIndex = `
	{
		"mappings": {
			"properties": {
				"tx_id": { "type": "keyword", "index": true },
				"tx_idx": { "type": "integer" },
				"idx": { "type": "integer" },
				"seq": { "type": "long" },
				"order": { "type": "long", "index": true },
				"close_time": { "type": "date" },
				"successful": { "type": "boolean" },
				"result_code": { "type": "byte" },
				"inner_result_code": { "type": "byte" },
				"tx_source_account_id": { "type": "keyword", "index": true },
				"memo": {
					"properties": {
						"type": { "type": "byte" },
						"value": { "type": "keyword" }
					}
				},
				"type": { "type": "keyword", "index": true },
				"source_account_id": { "type": "keyword", "index": true },
				"source_asset": {
					"properties": {
						"key": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"source_amount": { "type": "scaled_float", "scaling_factor": 10000000 },
				"destination_account_id": { "type": "keyword", "index": true },
				"destination_asset": {
					"properties": {
						"key": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"destination_amount": { "type": "scaled_float", "scaling_factor": 10000000 },
				"offer_price": { "type": "double" },
				"offer_price_n_d": {
					"properties": {
						"n": { "type": "integer" },
						"d": { "type": "integer" }
					}
				},
				"offer_id": { "type": "long" },
				"trust_limit": { "type": "scaled_float", "scaling_factor": 10000000 },
				"authorize": { "type": "boolean" },
				"bump_to": { "type": "long" },
				"path": {
					"properties": {
						"key": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"thresholds": {
					"properties": {
						"low": { "type": "byte" },
						"medium": { "type": "byte" },
						"high": { "type": "byte" },
						"master": { "type": "byte" }
					}
				},
				"home_domain": { "type": "keyword" },
				"inflation_dest_id": { "type": "keyword" },
				"set_flags": {
					"properties": {
						"required": { "type": "boolean" },
						"revocable": { "type": "boolean" },
						"immutable": { "type": "boolean" }
					}
				},
				"clear_flags": {
					"properties": {
						"required": { "type": "boolean" },
						"revocable": { "type": "boolean" },
						"immutable": { "type": "boolean" }
					}
				},
				"signer": {
					"properties": {
						"key": { "type": "keyword" },
						"weight": { "type": "byte" }
					}
				},
				"data": {
					"properties": {
						"name": { "type": "keyword" },
						"value": { "type": "keyword" }
					}
				},
				"result_source_account_balance": { "type": "scaled_float", "scaling_factor": 10000000 }
			}
		}
	}
`

const balanceIndex = `
	{
		"mappings": {
			"properties": {
				"account_id": { "type": "keyword", "index": true },
				"balance": { "type": "scaled_float", "scaling_factor": 10000000 },
				"source": { "type": "keyword" },
				"created_at": { "type": "date" },
				"asset": {
					"properties": {
						"key": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
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
	refreshIndex(balanceIndexName, balanceIndex)
}

func refreshIndex(name string, body string) {
	res, err := config.ES.Indices.Get([]string{name})

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
	res, err := config.ES.Indices.Delete([]string{index})
	fatalIfError(res, err)
}

func createIndex(index string, body string) {
	create := config.ES.Indices.Create

	res, err := create(
		index,
		create.WithBody(strings.NewReader(body)),
		create.WithIncludeTypeName(false),
	)
	fatalIfError(res, err)
}
