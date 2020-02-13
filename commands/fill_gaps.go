package commands

import (
	"bytes"
	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
	"log"
	"sort"
)

const batchSize = 100
const INDEX_RETRIES_COUNT = 25

type FillGapsCommand struct {
	ES es.EsAdapter
}

func (cmd *FillGapsCommand) Execute() {
	minSeq, maxSeq := cmd.ES.MinMaxSeq()
	log.Printf("Min seq is %d, max seq is %d\n", minSeq, maxSeq)

	var missing []int

	for i := minSeq; i < maxSeq; i += batchSize {
		to := i + batchSize
		if to > maxSeq {
			to = maxSeq
		}

		log.Println(i, to)
		seqs := cmd.ES.GetLedgerSeqsInRange(i, to)
		log.Println("Seqs ingested:", seqs)

		missing = append(missing, cmd.findMissing(i, to, seqs)...)
		log.Println("=============================")
	}

	if !*config.FillGapsDryRun {
		cmd.exportSeqs(missing)
	} else {
		log.Println(missing)
	}
}

func (cmd *FillGapsCommand) findMissing(from, to int, sortedArr []int) (missing []int) {
	rangeArr := makeRange(from, to)

	index := sort.SearchInts(sortedArr, from)

	if index == 0 && sortedArr[0] != from {
		index = sort.SearchInts(rangeArr, sortedArr[0])
		missing = append(missing, rangeArr[:index]...)
	}

	for i := index; i < len(rangeArr); i += 1 {
		if i-index >= len(sortedArr) || sortedArr[i-index] != rangeArr[i] {
			missing = append(missing, rangeArr[i])
		}
	}

	log.Println("Missing:", missing)
	return
}

func (cmd *FillGapsCommand) exportSeqs(seqs []int) {
	var b bytes.Buffer

	rows := db.LedgerHeaderRowFetchBySeqs(seqs)

	for n := 0; n < len(rows); n++ {
		txs := db.TxHistoryRowForSeq(rows[n].LedgerSeq)
		fees := db.TxFeeHistoryRowsForRows(txs)

		es.SerializeLedger(rows[n], txs, fees, &b)
	}

	cmd.ES.IndexWithRetries(&b, INDEX_RETRIES_COUNT)
}

func makeRange(from, to int) []int {
	a := make([]int, to-from)
	for i := range a {
		a[i] = from + i
	}
	return a
}
