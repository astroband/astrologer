package db

import (
	"bytes"
	"log"
	"net/url"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

// Copy paste from Horizon
func utf8Scrub(in string) string {

	// First check validity using the stdlib, returning if the string is already
	// valid
	if utf8.ValidString(in) {
		return in
	}

	left := []byte(in)
	var result bytes.Buffer

	for len(left) > 0 {
		r, n := utf8.DecodeRune(left)

		_, err := result.WriteRune(r)
		if err != nil {
			panic(err)
		}

		left = left[n:]
	}

	return result.String()
}

// Adapter defines the interface to work with ledger database
type Adapter interface {
	LedgerHeaderRowCount(first int, last int) int
	LedgerHeaderRowFetchBatch(n int, start int, batchSize int) []LedgerHeaderRow
	LedgerHeaderLastRow() *LedgerHeaderRow
	LedgerHeaderFirstRow() *LedgerHeaderRow
	LedgerHeaderNext(seq int) *LedgerHeaderRow
	LedgerHeaderGaps() (r []Gap)
	TxHistoryRowForSeq(seq int) []TxHistoryRow
	TxFeeHistoryRowsForRows(rows []TxHistoryRow) []TxFeeHistoryRow
}

// Client is an adapter implementation for stellar-core database
type Client struct {
	rawClient *sqlx.DB
}

// Connect returns the Client configured for the specified database
func Connect(databaseURL *url.URL) *Client {
	databaseDriver := (*databaseURL).Scheme

	db, err := sqlx.Connect(databaseDriver, (*databaseURL).String())
	if err != nil {
		log.Fatal(err)
	}

	return &Client{rawClient: db}
}
