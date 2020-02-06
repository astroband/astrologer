package db

import (
	"database/sql"
	"log"

	"github.com/astroband/astrologer/config"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/xdr"
)

// LedgerHeaderRowBatchSize used in LedgerHeaderRowFetchBatch
var LedgerHeaderRowBatchSize = *config.BatchSize

// LedgerHeaderRow is struct representing ledger in database
type LedgerHeaderRow struct {
	Hash           string           `db:"ledgerhash"`
	PrevHash       string           `db:"prevhash"`
	BucketListHash string           `db:"bucketlisthash"`
	LedgerSeq      int              `db:"ledgerseq"`
	CloseTime      int64            `db:"closetime"`
	Data           xdr.LedgerHeader `db:"data"`
}

// Gap represents gap in ledger sequence
type Gap struct {
	Start int `db:"gap_start"`
	End   int `db:"gap_end"`
}

// LedgerHeaderRowCount returns total ledgers count within given range
func LedgerHeaderRowCount(first int, last int) int {
	total := 0

	if last == 0 {
		config.DB.Get(&total, "SELECT count(ledgerseq) FROM ledgerheaders WHERE ledgerseq >= $1", first)
	} else {
		config.DB.Get(&total, "SELECT count(ledgerseq) FROM ledgerheaders WHERE ledgerseq >= $1 AND ledgerseq <= $2", first, last)
	}

	return total
}

// LedgerHeaderRowFetchBatch gets bunch of ledgers
func LedgerHeaderRowFetchBatch(n int, start int) []LedgerHeaderRow {
	ledgers := []LedgerHeaderRow{}
	offset := n * LedgerHeaderRowBatchSize
	low := offset + start
	high := low + LedgerHeaderRowBatchSize - 1

	err := config.DB.Select(
		&ledgers,
		"SELECT * FROM ledgerheaders WHERE ledgerseq BETWEEN $1 AND $2 ORDER BY ledgerseq ASC",
		low,
		high)

	if err != nil {
		log.Fatal(err)
	}

	return ledgers
}

func LedgerHeaderRowFetchBySeqs(seqs []int) []LedgerHeaderRow {
	ledgers := []LedgerHeaderRow{}

	query, _, err := sqlx.In("SELECT * FROM ledgerheaders WHERE ledgerseq IN (?) ORDER BY ledgerseq ASC;", seqs)

	if err != nil {
		log.Fatal(err)
	}

	query = config.DB.Rebind(query)
	err = config.DB.Select(&ledgers, query, ledgers)

	if err != nil {
		log.Fatal(err)
	}

	return ledgers
}

// LedgerHeaderLastRow returns lastest ledger in the database
func LedgerHeaderLastRow() *LedgerHeaderRow {
	var h LedgerHeaderRow

	err := config.DB.Get(&h, "SELECT * FROM ledgerheaders ORDER BY ledgerseq DESC LIMIT 1")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	return &h
}

// LedgerHeaderFirstRow returns lastest first ledger in the database
func LedgerHeaderFirstRow() *LedgerHeaderRow {
	var h LedgerHeaderRow

	err := config.DB.Get(&h, "SELECT * FROM ledgerheaders ORDER BY ledgerseq ASC LIMIT 1")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	return &h
}

// LedgerHeaderNext returns next ledger to fetch
func LedgerHeaderNext(seq int) *LedgerHeaderRow {
	var h LedgerHeaderRow

	err := config.DB.Get(&h, "SELECT * FROM ledgerheaders WHERE ledgerseq > $1 ORDER BY ledgerseq ASC LIMIT 1", seq)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	return &h
}

// LedgerHeaderGaps returns gap positions in ledgerheaders
func LedgerHeaderGaps() (r []Gap) {
	err := config.DB.Select(&r, `
		SELECT ledgerseq + 1 AS gap_start, next_nr - 1 AS gap_end
		FROM (
  		SELECT ledgerseq, LEAD(ledgerseq) OVER (ORDER BY ledgerseq) AS next_nr
  		FROM ledgerheaders
		) nr
		WHERE ledgerseq + 1 <> next_nr
	`)

	if err != nil {
		log.Fatal(err)
	}

	return r
}
