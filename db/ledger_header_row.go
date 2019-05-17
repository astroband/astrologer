package db

import (
	"database/sql"
	"log"

	"github.com/gzigzigzeo/stellar-core-export/config"
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

// LedgerHeaderRowCount returns total ledgers count
func LedgerHeaderRowCount(start int) int {
	count := 0
	config.DB.Get(&count, "SELECT count(ledgerseq) FROM ledgerheaders WHERE ledgerseq >= $1", start)
	return count
}

// LedgerHeaderRowFetchBatch gets bunch of ledgers
func LedgerHeaderRowFetchBatch(n int, start int) []LedgerHeaderRow {
	ledgers := []LedgerHeaderRow{}
	offset := n * LedgerHeaderRowBatchSize

	err := config.DB.Select(
		&ledgers,
		"SELECT * FROM ledgerheaders WHERE ledgerseq >= $1 ORDER BY ledgerseq ASC OFFSET $2 LIMIT $3",
		start,
		offset,
		LedgerHeaderRowBatchSize)

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
