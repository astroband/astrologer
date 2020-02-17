package commands

import (
	"fmt"
	"log"

	"github.com/astroband/astrologer/es"
)

type CreateIndexCommandConfig struct {
	Force bool
}

type CreateIndexCommand struct {
	ES     es.Adapter
	Config CreateIndexCommandConfig
}

// CreateIndex calls create-indexes command
func (cmd *CreateIndexCommand) Execute() {
	for name, def := range es.GetIndexDefinitions() {
		cmd.refreshIndex(name, def)
	}
	fmt.Println("Indicies created successfully!")
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
