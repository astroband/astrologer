# Installation

```
  go get git@github.com/astroband/stellar-core-export # not sure
  brew install elasticsearch
  ./stellar-core-export create-indexes
```

# Run

```
  ./stellar-core-export export
```

There are also --verbose and --dry-run flags for export.

# From scratch

```
  ./stellar-core-export create-indexes --force
  ./stellar-core-export export
```

# Postman

https://www.getpostman.com/downloads

See es.postman_collection.json

# Cluster storage size

```curl localhost:9200/_cluster/stats?human\&pretty | more```