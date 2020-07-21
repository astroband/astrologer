package commands

import (
	"fmt"
	"log"

	"github.com/astroband/astrologer/es"
)

// CreateIndexCommandConfig represents the configuration options for the `create-index` command
type CreateIndexCommandConfig struct {
	Force bool
}

// CreateIndexCommand represents the `create-index` CLI command
type CreateIndexCommand struct {
	ES     es.Adapter
	Config CreateIndexCommandConfig
}

// Execute creates Astrologer indices in ElasticSearch
func (cmd *CreateIndexCommand) Execute() {
	for name, def := range es.GetIndexDefinitions() {
		cmd.refreshIndex(name, def)
	}
	fmt.Println("Indices were created successfully!")
}

func (cmd *CreateIndexCommand) refreshIndex(name es.IndexName, schema es.IndexDefinition) {
	if !cmd.ES.IndexExists(name) {
		cmd.ES.CreateIndex(name, schema)
		log.Printf("%s index created!", name)
	} else {
		if cmd.Config.Force {
			cmd.ES.DeleteIndex(name)
			cmd.ES.CreateIndex(name, schema)
			log.Printf("%s index recreated!", name)
		} else {
			log.Printf("%s index found, skipping...", name)
		}
	}
}
