package main

import (
	"bytes"
	"log"
	"time"

	"github.com/astroband/astrologer/commands"
	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
	"github.com/gammazero/workerpool"
)

var (
	pool = workerpool.New(*config.Concurrency)
)

func main() {
	switch config.Command {
	case "stats":
		commands.Stats()
	case "create-index":
		commands.CreateIndex()
	case "export":
		commands.Export()
	case "ingest":
		ingest()
	}
}

func index(b *bytes.Buffer, n int) {
	res, err := config.ES.Bulk(bytes.NewReader(b.Bytes()))

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil || res.IsError() {
		if n > 5 {
			log.Fatal("5 retries for bulk failed, aborting")
		}

		time.Sleep(15 * time.Second)

		index(b, n+1)
	}
}

func ingest() {
	var h *db.LedgerHeaderRow

	if *config.StartIngest == 0 {
		h = db.LedgerHeaderLastRow()
	} else {
		if *config.StartIngest > 0 {
			h = db.LedgerHeaderNext(*config.StartIngest)
		} else {
			last := db.LedgerHeaderLastRow()

			if last == nil {
				log.Fatal("Nothing to ingest")
			}

			h = db.LedgerHeaderNext(last.LedgerSeq + *config.StartIngest)
		}
	}

	if h == nil {
		log.Fatal("Nothing to ingest")
	}

	log.Println("Starting ingest from", h.LedgerSeq)

	for {
		var b bytes.Buffer
		var seq = h.LedgerSeq

		txs := db.TxHistoryRowForSeq(seq)
		fees := db.TxFeeHistoryRowsForRows(txs)
		es.MakeBulk(*h, txs, fees, &b)

		pool.Submit(func() { index(&b, 0) })

		log.Println("Ledger", seq, "ingested.")

		h = db.LedgerHeaderNext(seq)

		for {
			if h != nil {
				break
			}
			time.Sleep(1 * time.Second)
			h = db.LedgerHeaderNext(seq)
		}
	}
}
