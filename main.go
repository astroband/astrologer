package main

import (
	cmd "github.com/astroband/astrologer/commands"
	cfg "github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
	"github.com/stellar/go/network"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.Version(cfg.Version)
	commandName := kingpin.Parse()

	esClient := es.Connect((*cfg.EsURL).String())

	var command cmd.Command

	switch commandName {
	case "stats":
		dbClient := db.Connect(*cfg.DatabaseURL)
		command = &cmd.StatsCommand{ES: esClient, DB: dbClient}
	case "create-index":
		config := cmd.CreateIndexCommandConfig{Force: *cfg.ForceRecreateIndexes}
		command = &cmd.CreateIndexCommand{ES: esClient, Config: config}
	case "export":
		var networkPassphrase string

		switch *cfg.Network {
		case "public":
			networkPassphrase = network.PublicNetworkPassphrase
		case "test":
			networkPassphrase = network.TestNetworkPassphrase
		}

		config := cmd.ExportCommandConfig{
			Start:             *cfg.Start,
			Count:             *cfg.Count,
			DryRun:            *cfg.ExportDryRun,
			RetryCount:        *cfg.Retries,
			BatchSize:         *cfg.BatchSize,
			NetworkPassphrase: networkPassphrase,
		}
		command = &cmd.ExportCommand{ES: esClient, Config: config}
	case "ingest":
		dbClient := db.Connect(*cfg.DatabaseURL)
		command = &cmd.IngestCommand{ES: esClient, DB: dbClient}
	case "es-stats":
		command = &cmd.EsStatsCommand{ES: esClient}
	}

	command.Execute()
}
