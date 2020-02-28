package commands

import (
	"github.com/astroband/astrologer/config"
	"github.com/gammazero/workerpool"
)

var (
	pool = workerpool.New(*config.Concurrency)
)

// Command is an interface representing an Astrologer CLI command
type Command interface {
	Execute()
}
