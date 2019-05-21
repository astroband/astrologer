package db

import (
	"log"

	"github.com/astroband/astrologer/config"
	"github.com/jmoiron/sqlx"
	"github.com/stellar/go/xdr"
)

// TxFeeHistoryRow represents row of txhistory table
type TxFeeHistoryRow struct {
	TxID      string                 `db:"txid"`
	LedgerSeq int                    `db:"ledgerseq"`
	Index     int                    `db:"txindex"`
	Changes   xdr.LedgerEntryChanges `db:"txchanges"`
}

// TxFeeHistoryRowsForRows returns transactions for specified ledger sorted by index
func TxFeeHistoryRowsForRows(rows []TxHistoryRow) []TxFeeHistoryRow {
	txs := []TxFeeHistoryRow{}

	if len(rows) == 0 {
		return txs
	}

	ids := make([]int, len(rows))

	for n := 0; n < len(rows); n++ {
		ids[n] = rows[n].LedgerSeq
	}

	query, args, err := sqlx.In("SELECT * FROM txfeehistory WHERE ledgerseq IN (?) ORDER BY ledgerseq, txindex", ids)
	if err != nil {
		log.Fatal(err)
	}

	query = config.DB.Rebind(query)
	err = config.DB.Select(&txs, query, args...)
	if err != nil {
		log.Fatal(err)
	}

	return txs
}
