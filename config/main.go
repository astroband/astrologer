package config

import (
	"fmt"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Version Application version
const Version string = "0.2.0"

// NumberWithSign represents arg having sign
type NumberWithSign struct {
	Value    int
	Explicit bool // True if + or - was passed
}

func (n *NumberWithSign) Set(value string) error {
	v, err := strconv.Atoi(value)

	if err != nil {
		return err
	}

	n.Value = v
	n.Explicit = value[0] == '-' || value[0] == '+'

	return nil
}

func (n *NumberWithSign) String() string {
	if n.Value == 0 {
		return "0"
	}

	if n.Explicit {
		return fmt.Sprintf("%+d", n.Value)
	} else {
		return fmt.Sprintf("%d", n.Value)
	}
}

// Helper for the kingpin argument custom parser
// See https://github.com/alecthomas/kingpin#custom-parsers
func NumberWithSignParse(s kingpin.Settings) (target *NumberWithSign) {
	target = &NumberWithSign{0, false}
	s.SetValue(target)
	return
}

var (
	createIndexCommand = kingpin.Command("create-index", "Create ES indexes")
	exportCommand      = kingpin.Command("export", "Run export")
	ingestCommand      = kingpin.Command("ingest", "Start real time ingestion")
	fastReplayCommand  = kingpin.Command("fast-replay", "Experiment with using stellar-core fast in-memory replay catchup")
	_                  = kingpin.Command("stats", "Print database ledger statistics")
	_                  = kingpin.Command("es-stats", "Print ES ranges stats")

	// DatabaseURL Stellar Core database URL
	DatabaseURL = kingpin.
			Flag("database-url", "Stellar Core database URL").
			Default("postgres://localhost/core?sslmode=disable").
			OverrideDefaultFromEnvar("DATABASE_URL").
			URL()

	// EsURL ElasticSearch URL
	EsURL = kingpin.
		Flag("es-url", "ElasticSearch URL").
		Default("http://localhost:9200").
		OverrideDefaultFromEnvar("ES_URL").
		URL()

	// Concurrency How many tasks and goroutines to produce (all at once for now)
	Concurrency = kingpin.
			Flag("concurrency", "Concurrency for indexing").
			Short('c').
			Default("5").
			OverrideDefaultFromEnvar("CONCURRENCY").
			Int()

	// BatchSize Batch size for bulk export
	BatchSize = exportCommand.
			Flag("batch", "Ledger batch size").
			Short('b').
			Default("50").
			Int()

	// Retries Number of retries
	Retries = exportCommand.
		Flag("retries", "Retries count").
		Default("25").
		Int()

	Start = exportCommand.Arg("start", "Ledger to start indexing, +100 means offset 100 from the first").Default("0").Int()

	// Count ledgers
	Count             = exportCommand.Arg("count", "Count of ledgers to ingest, should be aliquout batch size").Default("0").Int()
	NetworkPassphrase = exportCommand.
				Flag("network-passphrase", "Network passphrase to use").
				Default("Test SDF Network ; September 2015").
				OverrideDefaultFromEnvar("NETWORK_PASSPHRASE").
				String()

	// StartIngest ledger to start with ingesting
	StartIngest = ingestCommand.Arg("start", "Ledger to start ingesting").Int()

	// Verbose print data
	Verbose = exportCommand.Flag("verbose", "Print indexed data").Bool()

	// ExportDryRun do not index data
	ExportDryRun = exportCommand.Flag("dry-run", "Do not send actual data to Elastic").Bool()

	// ForceRecreateIndexes Allows indexes to be deleted before creation
	ForceRecreateIndexes = createIndexCommand.Flag("force", "Delete indexes before creation").Bool()

	FastReplayUpTo  = fastReplayCommand.Arg("upto", "Ledger to start indexing").Int()
	FastReplayCount = fastReplayCommand.Arg("count", "Ledgers count to catchup").Int()
)
