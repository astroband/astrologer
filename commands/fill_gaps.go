package commands

import (
	"bytes"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
	"github.com/astroband/astrologer/support"
	"log"
)

const batchSize = 100
const INDEX_RETRIES_COUNT = 25

type FillGapsCommandConfig struct {
	DryRun bool
	Start  *int
	Count  *int
}

type FillGapsCommand struct {
	ES     es.Adapter
	DB     db.Adapter
	Config *FillGapsCommandConfig

	minSeq int
	maxSeq int
	count  int
}

func (cmd *FillGapsCommand) Execute() {
	if cmd.Config.Start != nil && cmd.Config.Count != nil {
		cmd.minSeq = *cmd.Config.Start
		cmd.count = *cmd.Config.Count
		cmd.maxSeq = cmd.minSeq + cmd.count
	} else {
		var empty bool
		cmd.minSeq, cmd.maxSeq, empty = cmd.ES.MinMaxSeq()

		if empty {
			log.Println("ES is empty")
			return
		}

		if cmd.Config.Start != nil {
			cmd.minSeq = *cmd.Config.Start
			cmd.count = cmd.maxSeq - cmd.minSeq + 1
		}
	}

	log.Printf("Min seq is %d, max seq is %d\n", cmd.minSeq, cmd.maxSeq)

	var missing []int

	for i := cmd.minSeq; i < cmd.maxSeq; i += batchSize {
		to := i + batchSize
		if to > cmd.maxSeq {
			to = cmd.maxSeq
		}

		// log.Println(i, to)
		seqs := cmd.ES.GetLedgerSeqsInRange(i, to)
		// log.Println("Seqs ingested:", seqs)

		if len(seqs) > 0 {
			missing = append(missing, cmd.findMissing(seqs)...)
		} else {
			missing = append(missing, support.MakeRange(i-1, to)...)
		}
		// log.Println("=============================")
	}

	log.Println(missing)
	log.Println(len(missing))

	if !cmd.Config.DryRun {
		cmd.exportSeqs(missing)
	}
}

func (cmd *FillGapsCommand) findMissing(sortedArr []int) (missing []int) {
	for i := 1; i < len(sortedArr); i += 1 {
		diff := sortedArr[i] - sortedArr[i-1]
		if diff > 1 {
			missing = append(missing, support.MakeRange(sortedArr[i-1], sortedArr[i])...)
		}
	}

	// log.Println("Missing:", missing)
	return
}

func (cmd *FillGapsCommand) exportSeqs(seqs []int) {
	// exportConfig := ExportCommandConfig{
	// 	Start: config.NumberWithSign{Value: cmd.minSeq, Explicit: false},
	// 	Count: cmd.count,
	// }

	// exportCommand = &cmd.ExportCommand{ES: cmd.ES, DB: cmd.DB, Config: exportConfig}
	log.Printf("Exporting %d ledgers\n", len(seqs))

	var dbSeqs []int
	var b bytes.Buffer
	batch := 50

	for i := 0; i < len(seqs); i += batch {
		b.Reset()

		to := i + batch
		if to > len(seqs) {
			to = len(seqs) - 1
		}

		var seqsBlock []int

		if len(seqs) == 1 {
			seqsBlock = seqs
		} else {
			seqsBlock = seqs[i:to]
		}

		pool.Submit(func() {
			rows := cmd.DB.LedgerHeaderRowFetchBySeqs(seqsBlock)

			for n := 0; n < len(rows); n++ {
				log.Printf("Ingesting %d ledger\n", rows[n].LedgerSeq)
				dbSeqs = append(dbSeqs, rows[n].LedgerSeq)

				txs := cmd.DB.TxHistoryRowForSeq(rows[n].LedgerSeq)
				fees := cmd.DB.TxFeeHistoryRowsForRows(txs)

				es.SerializeLedger(rows[n], txs, fees, &b)
			}

			log.Println("Calling bulk insert")
			cmd.ES.IndexWithRetries(&b, INDEX_RETRIES_COUNT)
		})
	}

	pool.StopWait()
	diff := support.Difference(seqs, dbSeqs)

	if len(diff) > 0 {
		log.Printf("DB misses next ledgers: %v", diff)
	}
}
