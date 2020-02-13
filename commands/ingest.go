package commands

import (
	"bytes"
	"log"
	"time"

	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
)

const INGEST_RETRIES = 25

type IngestCommand struct {
	ES es.EsAdapter
}

// Execute Starts ingestion
func (cmd *IngestCommand) Execute() {
	current := cmd.getStartLedger()
	log.Println("Starting ingest from", current.LedgerSeq)

	for {
		var b bytes.Buffer
		var seq = current.LedgerSeq

		txs := db.TxHistoryRowForSeq(seq)
		fees := db.TxFeeHistoryRowsForRows(txs)

		es.SerializeLedger(*current, txs, fees, &b)
		//es.NewBulkMaker(*current, txs, fees, &b).Make()

		cmd.ES.IndexWithRetries(&b, INGEST_RETRIES)

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

func (cmd *IngestCommand) getStartLedger() (h *db.LedgerHeaderRow) {
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
