package es

// IndexName represents the name of ElasticSearch index
type IndexName string

// IndexDefinition represents the definition of ElasticSearch index
type IndexDefinition string

const (
	ledgerHeaderIndexName  IndexName = "ledger"
	txIndexName            IndexName = "tx"
	opIndexName            IndexName = "op"
	balanceIndexName       IndexName = "balance"
	tradesIndexName        IndexName = "trades"
	signerHistoryIndexName IndexName = "signers"
)

// GetIndexDefinitions returns ElasticSearch index definitions for Astrologer indices
func GetIndexDefinitions() map[IndexName]IndexDefinition {
	m := make(map[IndexName]IndexDefinition)

	m[ledgerHeaderIndexName] = `
      {
          "settings": {
            "index" : {
              "sort.field" : "paging_token",
              "sort.order" : "desc",
              "number_of_shards" : 1
            }
          },
          "mappings": {
            "properties": {
              "id": { "type": "keyword", "index": true },
              "hash": { "type": "keyword", "index": true },
              "prev_hash": { "type": "keyword", "index": false },
              "bucket_list_hash": { "type": "keyword", "index": false },
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

	m[txIndexName] = `
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
				"max_fee": { "type": "long" },
				"fee_charged": { "type": "long" },
        "fee_account_id": { "type": "keyword", "index": true },
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
				}
			}
		}
	}
`
	m[opIndexName] = `
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
        "authorize": { "type": "keyword" },
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
	m[balanceIndexName] = `
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

	m[tradesIndexName] = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 1
			}
		},
		"mappings": {
			"properties": {
        "id": { "type": "keyword", "index": true },
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

	m[signerHistoryIndexName] = `
	{
		"settings": {
			"index" : {
        "sort.field" : "paging_token",
				"sort.order" : "desc",
				"number_of_shards" : 1
			}
		},
		"mappings": {
			"properties": {
        "id": { "type": "keyword", "index": true },
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

	return m
}
