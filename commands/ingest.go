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
	DB db.DbAdapter
}

// Execute Starts ingestion
func (cmd *IngestCommand) Execute() {
	current := cmd.getStartLedger()
	log.Println("Starting ingest from", current.LedgerSeq)

	for {
		var b bytes.Buffer
		var seq = current.LedgerSeq

		txs := cmd.DB.TxHistoryRowForSeq(seq)
		fees := cmd.DB.TxFeeHistoryRowsForRows(txs)

		es.SerializeLedger(*current, txs, fees, &b)
		//es.NewBulkMaker(*current, txs, fees, &b).Make()

		cmd.ES.IndexWithRetries(&b, INGEST_RETRIES)

		log.Println("Ledger", seq, "ingested.")

		current = cmd.DB.LedgerHeaderNext(seq)

		for {
			if current != nil {
				break
			}
			time.Sleep(1 * time.Second)
			current = cmd.DB.LedgerHeaderNext(seq)
		}
	}
}

func (cmd *IngestCommand) getStartLedger() (h *db.LedgerHeaderRow) {
	if *config.StartIngest == 0 {
		h = cmd.DB.LedgerHeaderLastRow()
	} else {
		if *config.StartIngest > 0 {
			h = cmd.DB.LedgerHeaderNext(*config.StartIngest)
		} else {
			last := cmd.DB.LedgerHeaderLastRow()

			if last == nil {
				log.Fatal("Nothing to ingest")
			}

			h = cmd.DB.LedgerHeaderNext(last.LedgerSeq + *config.StartIngest)
		}
	}

	if h == nil {
		log.Fatal("Nothing to ingest")
	}

	return h
}
