package es

import (
	"context"
	"io"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gzigzigzeo/stellar-core-export/config"
)

// BulkIndex calls elasticsearch bulk indexing API
func BulkIndex(body io.Reader) {
	req := esapi.BulkRequest{
		Body:    body,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), config.ES)
	fatalIfError(res, err)

	defer res.Body.Close()
}
