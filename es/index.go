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
	defer res.Body.Close()

	checkErr(res, err)
}

// func index(index string, id string, body string) {
// 	req := esapi.IndexRequest{
// 		Index:      index,
// 		DocumentID: id,
// 		Body:       strings.NewReader(body),
// 		Refresh:    "true",
// 	}

// 	res, err := req.Do(context.Background(), config.ES)
// 	defer res.Body.Close()

// 	checkErr(res, err)
// }
