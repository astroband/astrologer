package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// SerializeForBulk returns object serialized for elastic bulk indexing
func SerializeForBulk(obj Indexable, b *bytes.Buffer) {
	var meta string

	id := obj.DocID()

	if id != nil {
		meta = fmt.Sprintf(
			`{ "create": { "_index": "%s", "_id": "%s", "_type": "_doc" } }%s`,
			obj.IndexName(),
			*id,
			"\n",
		)
	} else {
		meta = fmt.Sprintf(
			`{ "create": { "_index": "%s", "_type": "_doc" } }%s`, obj.IndexName(), "\n",
		)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}

	data = append(data, "\n"...)

	b.Grow(len(meta) + len(data))
	b.Write([]byte(meta))
	b.Write(data)
}
