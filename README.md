# Astrologer

Exports historical data to ElasticSearch storage (some data are still WIP).

# Installation

```
  go get git@github.com/astroband/astrologer
```

# Creating indexes

```
  ./astrologer create-indexes
```

Use `--force` flag to force recreate from scratch.

# Export from scratch

```
  ./astrologer export
```

You may use starting ledger number as second argument and ledger count as third. Note that real ledger count will be related to `--batch` parameter value, eg. if you specify start 0, count 150 and batch 100, 200 ledgers will be exported.

There are also `--verbose` and `--dry-run` flags for debug purposes.

# Ingest

```
  ./astrologer ingest
```

Will start live ingest from lastest ledger. You may use starting ledger number as second argument or specify some starting ledger in near past (useful for deployment):

```
  ./astrologer ingest -- -100
```

Will start ingestion from current ledger -100

# Postman

There are some example queries (aggregations mostly) in PostMan format.

https://www.getpostman.com/downloads

See `es.postman_collection.json`

# Check cluster storage size

```curl localhost:9200/_cluster/stats?human\&pretty | more```