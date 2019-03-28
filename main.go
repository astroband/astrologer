package main

import (
	"log"
	"strings"

	"github.com/gzigzigzeo/stellar-core-export/config"
	"github.com/gzigzigzeo/stellar-core-export/db"
	"github.com/gzigzigzeo/stellar-core-export/es"
	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	switch config.Command {
	case "create-index":
		es.CreateIndicies()
		log.Println("Indicies created successfully!")
	case "export":
		export()
	}
}

func export() {
	count := db.LedgerHeaderRowCount()
	bar := pb.StartNew(count)

	blocks := count / db.LedgerHeaderRowBatchSize
	if count%db.LedgerHeaderRowBatchSize > 0 {
		blocks = blocks + 1
	}

	for i := 0; i < blocks; i++ {
		var b strings.Builder

		rows := db.LedgerHeaderRowFetchBatch(i)

		for n := 0; n < len(rows); n++ {
			h := es.NewLedgerHeader(&rows[n])
			es.SerializeForBulk(h, &b)

			txs := db.TxHistoryRowForSeq(h.Seq)
			for t := 0; t < len(txs); t++ {
				tx := es.NewTransaction(&txs[t], h.CloseTime)
				es.SerializeForBulk(tx, &b)
			}

			bar.Increment()
		}

		es.BulkIndex(strings.NewReader(b.String()))
	}

	bar.Finish()
}
