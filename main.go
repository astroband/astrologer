package main

import (
	"github.com/astroband/astrologer/commands"
	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/es"
	"github.com/gammazero/workerpool"
)

var (
	pool = workerpool.New(*config.Concurrency)
)

func main() {
	switch config.Command {
	case "stats":
		commands.Stats()
	case "create-index":
		commands.CreateIndex()
	case "export":
		commands.Export()
	case "ingest":
		commands.Ingest()
	case "es-stats":
		commands.EsStats()
	case "fill-gaps":
		commands.FillGaps(es.Adapter)
	}
}
