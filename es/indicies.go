package es

import (
	"log"
	"net/http"
	"strings"

	"github.com/astroband/astrologer/config"
)

const ledgerHeaderIndex = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 4
			}
		},
		"mappings": {
			"properties": {
				"hash": { "type": "keyword", "index": true },
				"prev_hash": { "type": "keyword", "index": false },
				"bucket_list_hash": { "type": "keyword", "index": false },
				"tx_set_result_hash": { "type": "keyword", "index": false },
				"seq": { "type": "long" },
				"paging_token": { "type": "keyword", "index": true },
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
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 4
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword", "index": true },
				"idx": { "type": "integer" },
				"seq": { "type": "long" },
				"paging_token": { "type": "keyword", "index": true },
				"fee": { "type": "long" },
				"fee_charged": { "type": "long" },
				"operation_count": { "type": "integer" },
				"close_time": { "type": "date" },
				"successful": { "type": "boolean" },
				"result_code": { "type": "integer" },
				"source_account_id": { "type": "keyword", "index": true },
				"time_bounds": {
					"properties": {
						"min_time": { "type": "long" },
						"max_time": { "type": "long" }
					}
				},
				"memo": {
					"properties": {
						"type": { "type": "byte" },
						"value": { "type": "keyword" }
					}
				},
        "meta": { "type": "keyword", "index": false },
        "fee_meta": { "type": "keyword", "index": false }
			}
		}
	}
`

const opIndex = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 4
			}
		},
		"mappings": {
			"properties": {
				"tx_id": { "type": "keyword", "index": true },
				"tx_idx": { "type": "integer" },
				"idx": { "type": "integer" },
				"seq": { "type": "long" },
				"paging_token": { "type": "keyword", "index": true },
				"close_time": { "type": "date" },
				"successful": { "type": "boolean" },
				"result_code": { "type": "integer" },
				"inner_result_code": { "type": "integer" },
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
						"id": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"source_amount": { "type": "scaled_float", "scaling_factor": 10000000 },
				"destination_account_id": { "type": "keyword", "index": true },
				"destination_asset": {
					"properties": {
						"id": { "type": "keyword" },
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
						"id": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"thresholds": {
					"properties": {
						"low": { "type": "integer" },
						"medium": { "type": "integer" },
						"high": { "type": "integer" },
						"master": { "type": "integer" }
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
						"id": { "type": "keyword" },
						"weight": { "type": "integer" },
						"type": { "type": "byte" }
					}
				},
				"data": {
					"properties": {
						"name": { "type": "keyword" },
						"value": { "type": "keyword" }
					}
				},
				"result_source_account_balance": { "type": "scaled_float", "scaling_factor": 10000000 },
				"result_offer": {
					"properties": {
						"amount": { "type": "scaled_float", "scaling_factor": 10000000 },
						"price": { "type": "scaled_float", "scaling_factor": 10000000 },
						"price_n_d": {
							"properties": {
								"n": { "type": "integer" },
								"d": { "type": "integer" }
							}
						},
						"selling": {
							"properties": {
								"id": { "type": "keyword" },
								"code": { "type": "keyword" },
								"issuer": { "type": "keyword" }
							}
						},
						"buying": {
							"properties": {
								"id": { "type": "keyword" },
								"code": { "type": "keyword" },
								"issuer": { "type": "keyword" }
							}
						},
						"offer_id": { "type": "long" },
						"seller_id": { "type": "keyword" }
					}
				},
				"result_offer_effect": { "type": "keyword" }
			}
		}
	}
`

const balanceIndex = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 4
			}
		},
		"mappings": {
			"properties": {
				"paging_token": { "type": "keyword", "index": true },
				"account_id": { "type": "keyword", "index": true },
				"value": { "type": "scaled_float", "scaling_factor": 10000000 },
				"diff": { "type": "scaled_float", "scaling_factor": 10000000 },
				"positive": { "type": "boolean", "index": true },
				"source": { "type": "keyword" },
				"created_at": { "type": "date" },
				"asset": {
					"properties": {
						"id": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				}
			}
		}
	}
`

const eventsIndex = `
{
  "settings": {
    "index" : {
      "sort.field" : "paging_token",
      "sort.order" : "desc",
      "number_of_shards" : 4
    }
  },
  "mappings": {
    "properties": {
      "paging_token": { "type": "keyword", "index": true },
      "account_id": { "type": "keyword", "index": true },
      "source": { "type": "keyword" },
      "created_at": { "type": "date" },
      "entity": { "type": "keyword", "index": true },
      "type": { "type": "keyword", "index": true },

      "offer_id": { "type": "long" },
      "seller": { "type": "keyword" },
      "selling": {
        "properties": {
          "id": { "type": "keyword" },
          "code": { "type": "keyword" },
          "issuer": { "type": "keyword" }
        }
      },
      "buying": {
        "properties": {
          "id": { "type": "keyword" },
          "code": { "type": "keyword" },
          "issuer": { "type": "keyword" }
        }
      },
      "amount": { "type": "scaled_float", "scaling_factor": 10000000 },
      "price": { "type": "scaled_float", "scaling_factor": 10000000 },
      "passive": { "type": "boolean" }
    }
  }
}
`

const tradesIndex = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 4
			}
		},
		"mappings": {
			"properties": {
				"paging_token": { "type": "keyword", "index": true },				
				"sold": { "type": "scaled_float", "scaling_factor": 10000000 },
				"bought": { "type": "scaled_float", "scaling_factor": 10000000 },
				"asset_sold": {
					"properties": {
						"id": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"asset_bought": {
					"properties": {
						"id": { "type": "keyword" },
						"code": { "type": "keyword" },
						"issuer": { "type": "keyword" }
					}
				},
				"offer_id": { "type": "long" },
				"seller_id": { "type": "keyword", "index": true },
				"buyer_id": { "type": "keyword", "index": true },
				"price": { "type": "scaled_float", "scaling_factor": 10000000 },
				"time": { "type": "date" }
			}
		}
	}
`

const signerHistoryIndex = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 4
			}
		},
		"mappings": {
			"properties": {
				"paging_token": { "type": "keyword", "index": true },
				"account_id": { "type": "keyword", "index": true },
				"signer": { "type": "keyword", "index": true },
				"type": { "type": "byte" },
				"weight": { "type": "integer" },
				"seq": { "type": "integer" },
				"tx_idx": { "type": "integer" },
				"idx": { "type": "integer" },
				"ledger_close_time": { "type": "date" }
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
	refreshIndex(eventIndexName, eventsIndex)
	refreshIndex(tradesIndexName, tradesIndex)
	refreshIndex(signerHistoryIndexName, signerHistoryIndex)
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
