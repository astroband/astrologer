package commands

import (
	"bytes"
	"log"
	"time"

	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
)

var (
	current = getStartLedger()
)

// Ingest Starts ingestion
func Ingest() {
	log.Println("Starting ingest from", current.LedgerSeq)

	for {
		var b bytes.Buffer
		var seq = current.LedgerSeq

		txs := db.TxHistoryRowForSeq(seq)
		fees := db.TxFeeHistoryRowsForRows(txs)

		es.SerializeLedger(*current, txs, fees, &b)
		//es.NewBulkMaker(*current, txs, fees, &b).Make()

		index(&b, 0) // Defined in export.go

		log.Println("Ledger", seq, "ingested.")

		current = db.LedgerHeaderNext(seq)

		for {
			if current != nil {
				break
			}
			time.Sleep(1 * time.Second)
			current = db.LedgerHeaderNext(seq)
		}
	}
}

func getStartLedger() (h *db.LedgerHeaderRow) {
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

	return h
}
