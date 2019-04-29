package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/gzigzigzeo/stellar-core-export/config"
	"github.com/gzigzigzeo/stellar-core-export/db"
	"github.com/gzigzigzeo/stellar-core-export/es"
	"github.com/ti/nasync"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	async = nasync.New(10, 10)
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

func index(b *bytes.Buffer, n int) {
	res, err := config.ES.Bulk(bytes.NewReader(b.Bytes()))
	defer res.Body.Close()

	if err != nil || res.IsError() {
		if n > 5 {
			log.Fatal("5 retries for bulk failed, aborting")
		}

		log.Println("Retrying...")
		time.Sleep(10 * time.Second)
		async.Do(index, b, n+1)
	}
}

func export() {
	count := db.LedgerHeaderRowCount()
	bar := pb.StartNew(count)

	blocks := count / db.LedgerHeaderRowBatchSize
	if count%db.LedgerHeaderRowBatchSize > 0 {
		blocks = blocks + 1
	}

	defer async.Close()

	for i := 0; i < blocks; i++ {
		var b bytes.Buffer

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
					id := fmt.Sprintf("%v:%v:%v", h.Seq, t, o)
					bl := es.ExtractBalances(metas[o].Changes, h.CloseTime, id)
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
			async.Do(index, &b, 0)
		}
	}

	if !*config.Verbose {
		bar.Finish()
	}
}
