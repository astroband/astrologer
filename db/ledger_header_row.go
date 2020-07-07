package db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/xdr"
	"log"
)

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
func (db *Client) LedgerHeaderRowCount(first, last int) int {
	total := 0

	if last == 0 {
		db.rawClient.Get(&total, "SELECT count(ledgerseq) FROM ledgerheaders WHERE ledgerseq >= $1", first)
	} else {
		db.rawClient.Get(&total, "SELECT count(ledgerseq) FROM ledgerheaders WHERE ledgerseq >= $1 AND ledgerseq <= $2", first, last)
	}

	return total
}

// LedgerHeaderRowFetchBatch gets bunch of ledgers
func (db *Client) LedgerHeaderRowFetchBatch(n, start, batchSize int) []LedgerHeaderRow {
	ledgers := []LedgerHeaderRow{}
	offset := n * batchSize
	low := offset + start
	high := low + batchSize - 1

	err := db.rawClient.Select(
		&ledgers,
		"SELECT * FROM ledgerheaders WHERE ledgerseq BETWEEN $1 AND $2 ORDER BY ledgerseq ASC",
		low,
		high)

	if err != nil {
		log.Fatal(err)
	}

	return ledgers
}

// LedgerHeaderLastRow returns lastest ledger in the database
func (db *Client) LedgerHeaderLastRow() *LedgerHeaderRow {
	var h LedgerHeaderRow

	err := db.rawClient.Get(&h, "SELECT * FROM ledgerheaders ORDER BY ledgerseq DESC LIMIT 1")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	return &h
}

// LedgerHeaderFirstRow returns lastest first ledger in the database
func (db *Client) LedgerHeaderFirstRow() *LedgerHeaderRow {
	var h LedgerHeaderRow

	err := db.rawClient.Get(&h, "SELECT * FROM ledgerheaders ORDER BY ledgerseq ASC LIMIT 1")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	return &h
}

// LedgerHeaderNext returns next ledger to fetch
func (db *Client) LedgerHeaderNext(seq int) *LedgerHeaderRow {
	var h LedgerHeaderRow

	err := db.rawClient.Get(&h, "SELECT * FROM ledgerheaders WHERE ledgerseq > $1 ORDER BY ledgerseq ASC LIMIT 1", seq)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	return &h
}

func (db *Client) LedgerHeaderRowFetchBySeqs(seqs []int) []LedgerHeaderRow {
	ledgers := []LedgerHeaderRow{}

	query, args, err := sqlx.In("SELECT * FROM ledgerheaders WHERE ledgerseq IN (?) ORDER BY ledgerseq ASC;", seqs)

	if err != nil {
		log.Fatal(err)
	}

	query = db.rawClient.Rebind(query)
	err = db.rawClient.Select(&ledgers, query, args...)

	if err != nil {
		log.Fatal(err)
	}

	return ledgers
}

// LedgerHeaderGaps returns gap positions in ledgerheaders
func (db *Client) LedgerHeaderGaps() (r []Gap) {
	err := db.rawClient.Select(&r, `
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
