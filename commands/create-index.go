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
	for _, index := range es.GetIndexDefinitions() {
		cmd.refreshIndex(index.Name, index.Schema)
	}
	fmt.Println("Indicies created successfully!")
}

func (cmd *CreateIndexCommand) refreshIndex(name string, schema string) {
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
