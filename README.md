[![CircleCI](https://img.shields.io/circleci/build/gh/astroband/astrologer/master)](https://circleci.com/gh/astroband/astrologer) [![Go Report Card](https://goreportcard.com/badge/astroband/astrologer)](https://goreportcard.com/report/astroband/astrologer) [![License](https://img.shields.io/github/license/astroband/astrologer)]((https://github.com/astroband/astrologer/blob/master/LICENSE)[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fastroband%2Fastrologer.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fastroband%2Fastrologer?ref=badge_shield)
)


# Astrologer

Exports historical data to ElasticSearch storage (some data are still WIP).

# Installation

```
  go get git@github.com/astroband/astrologer
```

# Creating indexes

```
  ./astrologer create-index
```

Use `--force` flag to force recreate from scratch.

# Export from scratch

```
  ./astrologer export                 # Export everything
  ./astrologer export 23269090        # Start at ledger 23269090
  ./astrologer export 23269090 100    # 100 ledgers starting with 23269090
  ./astrologer export +1000           # Skip first 1000 ledgers
  ./astrologer export +1000 1000      # Skip first 1000 ledgers, limit to 1000
  ./astrologer export -- -1000        # Last 1000 ledgers
  ./astrologer export -- -1000 500    # 500 ledgers, offset -1000 from last
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

# Stats

Reports ledger segments existing in database.

```
  ./astrologer stats

  +----------+----------+--------+
  |   MIN    |   MAX    | COUNT  |
  +----------+----------+--------+
  |       10 |       10 |      1 |
  |       22 |       22 |      1 |
  | 23268991 | 23368992 | 100002 |
  +----------+----------+--------+
  |             TOTAL   | 100004 |
  +----------+----------+--------+
```

# ES Stats

Reports ledger segments existing elastic database.

```
  ./astrologer es-stats
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fastroband%2Fastrologer.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fastroband%2Fastrologer?ref=badge_large)