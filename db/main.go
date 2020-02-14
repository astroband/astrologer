package db

import (
	"bytes"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
	"log"
	"net/url"
	"unicode/utf8"
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

type DbAdapter interface {
	LedgerHeaderRowCount(first int, last int) int
	LedgerHeaderRowFetchBatch(n int, start int, batchSize int) []LedgerHeaderRow
	LedgerHeaderRowFetchBySeqs(seqs []int) []LedgerHeaderRow
	LedgerHeaderLastRow() *LedgerHeaderRow
	LedgerHeaderFirstRow() *LedgerHeaderRow
	LedgerHeaderNext(seq int) *LedgerHeaderRow
	LedgerHeaderGaps() (r []Gap)
	TxHistoryRowForSeq(seq int) []TxHistoryRow
	TxFeeHistoryRowsForRows(rows []TxHistoryRow) []TxFeeHistoryRow
}

type DbClient struct {
	rawClient *sqlx.DB
}

func Connect(databaseUrl *url.URL) *DbClient {
	databaseDriver := (*databaseUrl).Scheme

	db, err := sqlx.Connect(databaseDriver, (*databaseUrl).String())
	if err != nil {
		log.Fatal(err)
	}

	return &DbClient{rawClient: db}
}
