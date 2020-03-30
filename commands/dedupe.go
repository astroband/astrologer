package commands

import (
	"github.com/astroband/astrologer/es"
  "log"
)

type DedupeCommandConfig struct {
	DryRun bool
}

type DedupeCommand struct {
	ES     es.Adapter
	Config DedupeCommandConfig
}

func (cmd *DedupeCommand) Execute() {
  for indexName, _ := range es.GetIndexDefinitions() {
    log.Printf("Searching '%s' index...\n", indexName)
    cmd.ES.FindDuplicates(string(indexName), "paging_token")
  }
}
