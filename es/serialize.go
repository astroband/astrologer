package es

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Serialize returns object serialized for elastic indexing
func Serialize(obj Indexable, b *strings.Builder) {
	enc := json.NewEncoder(b)
	err := enc.Encode(obj)
	if err != nil {
		log.Fatal(err)
	}
}

// SerializeForBuilk returns object serialized for elastic bulk indexing
func SerializeForBulk(obj Indexable, b *strings.Builder) {
	b.WriteString(fmt.Sprintf(
		`{ "index": { "_index": "%s", "_id": "%s", "_type": "_doc" } }`,
		obj.IndexName(),
		obj.DocID(),
	))

	b.WriteString("\n")
	Serialize(obj, b)
	b.WriteString("\n")
}
