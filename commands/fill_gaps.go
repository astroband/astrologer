package commands

import (
	"bufio"
	"bytes"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
	"github.com/astroband/astrologer/support"
	"log"
	"strings"
)

const batchSize = 1000
const INDEX_RETRIES_COUNT = 25

type FillGapsCommandConfig struct {
	DryRun    bool
	Start     *int
	Count     *int
	BatchSize *int
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
		var to int

		if i+batchSize > cmd.maxSeq {
			to = cmd.maxSeq
		} else {
			to = i + batchSize - 1
		}

		seqs := cmd.ES.GetLedgerSeqsInRange(i, to)

		if len(seqs) > 0 {
			missing = append(missing, cmd.findMissing(seqs)...)

			if seqs[len(seqs)-1] != to {
				missing = append(missing, support.MakeRangeGtLte(seqs[len(seqs)-1], to)...)
			}
		} else {
			missing = append(missing, support.MakeRangeGteLte(i, to)...)
		}
	}

	cmd.exportSeqs(missing)
}

func (cmd *FillGapsCommand) findMissing(sortedArr []int) (missing []int) {
	for i := 1; i < len(sortedArr); i += 1 {
		diff := sortedArr[i] - sortedArr[i-1]
		if diff > 1 {
			missing = append(missing, support.MakeRangeGtLt(sortedArr[i-1], sortedArr[i])...)
		}
	}

	// log.Println("Missing:", missing)
	return
}

func (cmd *FillGapsCommand) exportSeqs(seqs []int) {
	log.Printf("Exporting %d ledgers\n", len(seqs))

	var dbSeqs []int
	batchSize := *cmd.Config.BatchSize

	for i := 0; i < len(seqs); i += batchSize {

		to := i + batchSize
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
			var b bytes.Buffer
			rows := cmd.DB.LedgerHeaderRowFetchBySeqs(seqsBlock)

			for n := 0; n < len(rows); n++ {
				// log.Printf("Ingesting %d ledger\n", rows[n].LedgerSeq)
				dbSeqs = append(dbSeqs, rows[n].LedgerSeq)

				txs := cmd.DB.TxHistoryRowForSeq(rows[n].LedgerSeq)
				fees := cmd.DB.TxFeeHistoryRowsForRows(txs)

				es.SerializeLedger(rows[n], txs, fees, &b)
			}

			log.Printf(
				"Bulk inserting %d docs, total size is %s\n",
				countLines(b)/2,
				support.ByteCountBinary(b.Len()),
			)

			if !cmd.Config.DryRun {
				cmd.ES.BulkInsert(&b)
			}
		})
	}

	pool.StopWait()
	diff := support.Difference(seqs, dbSeqs)

	if len(diff) > 0 {
		log.Printf("DB misses next ledgers: %v", diff)
	}
}

func countLines(buf bytes.Buffer) int {
	scanner := bufio.NewScanner(strings.NewReader(buf.String()))

	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanLines)

	// Count the lines.
	count := 0
	for scanner.Scan() {
		count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("reading input:", err)
	}

	return count
}
