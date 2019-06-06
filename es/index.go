package es

import (
	"context"
	"io"

	"github.com/astroband/astrologer/config"
	"github.com/elastic/go-elasticsearch/esapi"
)

// BulkIndex calls elasticsearch bulk indexing API
func BulkIndex(body io.Reader) {
	req := esapi.BulkRequest{
		Body: body,
	}

	res, err := req.Do(context.Background(), config.ES)
	fatalIfError(res, err)

	defer res.Body.Close()
}
