package main

import (
	cmd "github.com/astroband/astrologer/commands"
	cfg "github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/es"
	"github.com/gammazero/workerpool"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	pool = workerpool.New(*cfg.Concurrency)
)

func main() {
	kingpin.Version(cfg.Version)
	commandName := kingpin.Parse()

	esClient := es.Connect((*cfg.EsUrl).String())

	var command cmd.Command

	switch commandName {
	case "stats":
		command = &cmd.StatsCommand{ES: esClient}
	case "create-index":
		config := cmd.CreateIndexCommandConfig{Force: *cfg.ForceRecreateIndexes}
		command = &cmd.CreateIndexCommand{ES: esClient, Config: config}
	case "export":
		config := cmd.ExportCommandConfig{
			Start:      cfg.Start,
			Count:      *cfg.Count,
			DryRun:     *cfg.ExportDryRun,
			RetryCount: *cfg.Retries,
		}
		command = &cmd.ExportCommand{ES: esClient, Config: config}
	case "ingest":
		command = &cmd.IngestCommand{ES: esClient}
	case "es-stats":
		command = &cmd.EsStatsCommand{ES: esClient}
	case "fill-gaps":
		command = &cmd.FillGapsCommand{ES: esClient}
	}

	command.Execute()
}
