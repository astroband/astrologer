package config

import (
	"log"

	es "github.com/elastic/go-elasticsearch"
	"github.com/jmoiron/sqlx"
	"gopkg.in/alecthomas/kingpin.v2"

	_ "github.com/lib/pq" // Postgres driver
)

// Version Application version
const Version string = "0.0.1"

var (
	createIndexCommand = kingpin.Command("create-index", "Create ES indexes")
	exportCommand      = kingpin.Command("export", "Run export")
	ingestCommand      = kingpin.Command("ingest", "Start real time ingestion")

	databaseURL = kingpin.
			Flag("database-url", "Stellar Core database URL").
			Default("postgres://localhost/core?sslmode=disable").
			OverrideDefaultFromEnvar("DATABASE_URL").
			URL()

	esURL = kingpin.
		Flag("es-url", "ElasticSearch URL").
		Default("http://localhost:9200").
		OverrideDefaultFromEnvar("ES_URL").
		URL()

	// IndexConcurrency How many tasks and goroutines to produce (all at once for now)
	IndexConcurrency = kingpin.
				Flag("index-concurrency", "Concurrency for indexing").
				Short('c').
				Default("3").
				OverrideDefaultFromEnvar("INDEX_CONCURRENCY").
				Int()

	// FetchConcurrency How many tasks and goroutines to produce (all at once for now)
	FetchConcurrency = exportCommand.
				Flag("fetch-concurrency", "Concurrency for fetching").
				Short('f').
				Default("3").
				OverrideDefaultFromEnvar("FETCH_CONCURRENCY").
				Int()

	// BatchSize Batch size for bulk export
	BatchSize = exportCommand.
			Flag("batch", "Ledger batch size").
			Default("50").
			Int()

	// Start ledger to start with
	Start = exportCommand.Arg("start", "Ledger to start indexing").Default("0").Int()

	// StartIngest ledger to start with ingesting
	StartIngest = ingestCommand.Arg("start", "Ledger to start ingesting").Int()

	// Verbose print data
	Verbose = exportCommand.Flag("verbose", "Print indexed data").Bool()

	// DryRun do not index data
	DryRun = exportCommand.Flag("dry-run", "Do not send actual data to Elastic").Bool()

	// ForceRecreateIndexes Allows indexes to be deleted before creation
	ForceRecreateIndexes = createIndexCommand.Flag("force", "Delete indexes before creation").Bool()

	// DB Instance of sqlx.DB
	DB *sqlx.DB

	// ES ElasticSearch client instance
	ES *es.Client

	// Command KingPin command
	Command string
)

func initDB() {
	databaseDriver := (*databaseURL).Scheme

	db, err := sqlx.Connect(databaseDriver, (*databaseURL).String())
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func initES() {
	esCfg := es.Config{
		Addresses: []string{(*esURL).String()},
	}

	client, err := es.NewClient(esCfg)
	if err != nil {
		log.Fatal(err)
	}

	ES = client
}

func init() {
	kingpin.Version(Version)
	Command = kingpin.Parse()

	initDB()
	initES()
}
