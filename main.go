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
				txRow := &txs[t]
				ops := txRow.Envelope.Tx.Operations
				metas := txRow.Meta.V1.Operations // TODO: V1.Operations || MustOperations

				tx := es.NewTransaction(txRow, h.CloseTime)
				es.SerializeForBulk(tx, &b)

				for o := 0; o < len(ops); o++ {
					op := es.NewOperation(tx, &ops[o], byte(o))
					es.SerializeForBulk(op, &b)
				}

				for o := 0; o < len(metas); o++ {
					bl := es.ExtractBalances(metas[o].Changes)
					for _, balance := range bl {
						es.SerializeForBulk(balance, &b)
					}
				}
			}

			if !*config.Verbose {
				bar.Increment()
			}
		}

		if *config.Verbose {
			log.Println(b.String())
		}

		if !*config.DryRun {
			es.BulkIndex(strings.NewReader(b.String()))
		}
	}

	if !*config.Verbose {
		bar.Finish()
	}
}
