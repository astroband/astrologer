package main

import (
	"bytes"
	"log"
	"time"

	"github.com/gzigzigzeo/stellar-core-export/config"
	"github.com/gzigzigzeo/stellar-core-export/db"
	"github.com/gzigzigzeo/stellar-core-export/es"
	"github.com/ti/nasync"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	indexPool = nasync.New(*config.IndexConcurrency, *config.IndexConcurrency)
	fetchPool = nasync.New(*config.FetchConcurrency, *config.FetchConcurrency)
)

func main() {
	switch config.Command {
	case "create-index":
		es.CreateIndicies()
		log.Println("Indicies created successfully!")
	case "export":
		export()
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

		log.Println("Retrying ", n)
		time.Sleep(5 * time.Second)
		indexPool.Do(index, b, n+1)
	}
}

func fetch(i int, bar *pb.ProgressBar) {
	var b bytes.Buffer

	rows := db.LedgerHeaderRowFetchBatch(i, *config.Start)

	for n := 0; n < len(rows); n++ {
		txs := db.TxHistoryRowForSeq(rows[n].LedgerSeq)

		es.MakeBulk(rows[n], txs, &b)

		if !*config.Verbose {
			bar.Increment()
		}
	}

	if *config.Verbose {
		log.Println(b.String())
	}

	if !*config.DryRun {
		indexPool.Do(index, &b, 0)
	}
}

func export() {
	defer indexPool.Close()
	defer fetchPool.Close()

	count := db.LedgerHeaderRowCount(*config.Start)
	bar := pb.StartNew(count)

	blocks := count / db.LedgerHeaderRowBatchSize
	if count%db.LedgerHeaderRowBatchSize > 0 {
		blocks = blocks + 1
	}

	for i := 0; i < blocks; i++ {
		fetchPool.Do(fetch, i, bar)
	}

	if !*config.Verbose {
		bar.Finish()
	}
}

func ingest() {
	var h *db.LedgerHeaderRow

	defer indexPool.Close()

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
		es.MakeBulk(*h, txs, &b)

		indexPool.Do(index, &b, 0)

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
