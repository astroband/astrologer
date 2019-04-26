package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

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

func worker(id int, jobs <-chan *bytes.Buffer) {
	for b := range jobs {
		res, err := config.ES.Bulk(bytes.NewReader(b.Bytes()))

		if err != nil {
			log.Fatal("Error bulk", err)
		}

		if res.IsError() {
			if res.StatusCode == http.StatusTooManyRequests {

			} else {
				log.Fatal("Error bulk", res)
			}
		}
	}
}

func export() {
	count := db.LedgerHeaderRowCount()
	bar := pb.StartNew(count)

	blocks := count / db.LedgerHeaderRowBatchSize
	if count%db.LedgerHeaderRowBatchSize > 0 {
		blocks = blocks + 1
	}

	jobs := make(chan *bytes.Buffer, 100)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs)
	}

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
			jobs <- &b
		}
	}

	close(jobs)

	if !*config.Verbose {
		bar.Finish()
	}
}
