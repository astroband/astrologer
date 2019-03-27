package db

import (
	"log"

	"github.com/gzigzigzeo/stellar-core-export/config"
)

// TxHistoryRow represents row of txhistory table
type TxHistoryRow struct {
	TxID      string `db:"txid"`
	LedgerSeq int    `db:"ledgerseq"`
	TxIndex   int    `db:"txindex"`
	TxBody    string `db:"txbody"`
	TxResult  string `db:"txresult"`
	TxMeta    string `db:"txmeta"`
}

// TxHistoryRowForSeq returns transactions for specified ledger sorted by index
func TxHistoryRowForSeq(seq int) []TxHistoryRow {
	txs := []TxHistoryRow{}

	err := config.DB.Select(&txs, "SELECT * FROM txhistory WHERE ledgerseq = $1 ORDER BY txindex", seq)

	if err != nil {
		log.Fatal(err)
	}

	return txs
}
