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
	"github.com/schollz/progressbar"
)

var (
	pool = workerpool.New(*config.Concurrency)
)

func main() {
	switch config.Command {
	case "stats":
		commands.Stats()
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

		time.Sleep(15 * time.Second)

		index(b, n+1)
	}
}

func fetch(i int, bar *progressbar.ProgressBar) {
	var b bytes.Buffer

	rows := db.LedgerHeaderRowFetchBatch(i, *config.Start)

	for n := 0; n < len(rows); n++ {
		txs := db.TxHistoryRowForSeq(rows[n].LedgerSeq)
		fees := db.TxFeeHistoryRowsForRows(txs)

		es.MakeBulk(rows[n], txs, fees, &b)

		if !*config.Verbose {
			bar.Add(1)
		}
	}

	if *config.Verbose {
		log.Println(b.String())
	}

	if !*config.DryRun {
		index(&b, 0)
	}
}

func export() {
	count := db.LedgerHeaderRowCount(*config.Start, *config.Count)

	if count == 0 {
		log.Fatal("Nothing to export!")
	}

	bar := progressbar.NewOptions(
		count,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionThrottle(500*time.Millisecond),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWidth(100),
	)

	bar.RenderBlank()

	blocks := count / db.LedgerHeaderRowBatchSize
	if count%db.LedgerHeaderRowBatchSize > 0 {
		blocks = blocks + 1
	}

	for i := 0; i < blocks; i++ {
		i := i
		pool.Submit(func() { fetch(i, bar) })
	}

	pool.StopWait()

	if !*config.Verbose {
		bar.Finish()
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
