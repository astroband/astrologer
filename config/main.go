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

	// DB Instance of sqlx.DB
	DB *sqlx.DB

	// ES ElasticSearch client instance
	ES *es.Client

	// Command KingPin command
	Command string

	// ForceRecreateIndexes Allows indexes to be deleted before creation
	ForceRecreateIndexes = createIndexCommand.Flag("force", "Delete indexes before creation").Bool()

	databaseURL = kingpin.
			Flag("database-url", "Stellar Core database URL").
			Default("postgres://localhost/core?sslmode=disable").
			URL()

	esURL = kingpin.
		Flag("es-url", "ElasticSearch URL").
		Default("http://localhost:9200").
		URL()
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
