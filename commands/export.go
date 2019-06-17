package commands

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/schollz/progressbar"

	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
)

var (
	bar         *progressbar.ProgressBar
	first, last = getRange()
)

// Export command
func Export() {
	total := db.LedgerHeaderRowCount(first, last)

	if total == 0 {
		log.Fatal("Nothing to export within given range!", first, last)
	}

	fmt.Println("Exporting ledgers from", first, "to", last, "total", total)

	createBar(total)

	for i := 0; i < blockCount(total); i++ {
		i := i
		pool.Submit(func() { exportBlock(i) })
	}

	pool.StopWait()
	finishBar()
}

func exportBlock(i int) {
	var b bytes.Buffer

	rows := db.LedgerHeaderRowFetchBatch(i, first)

	for n := 0; n < len(rows); n++ {
		txs := db.TxHistoryRowForSeq(rows[n].LedgerSeq)
		fees := db.TxFeeHistoryRowsForRows(txs)

		//es.NewBulkMaker(rows[n], txs, fees, &b).Make()
		es.SerializeLedger(rows[n], txs, fees, &b)

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

func index(b *bytes.Buffer, retry int) {
	res, err := config.ES.Bulk(bytes.NewReader(b.Bytes()))

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil || (res != nil && res.IsError()) {
		if retry > 5 {
			log.Fatal("5 retries for bulk failed, aborting")
		}

		delay := time.Duration((rand.Intn(10) + 5))
		time.Sleep(delay * time.Second)

		index(b, retry+1)
	}
}

// Parses range of export command
func getRange() (first int, last int) {
	firstLedger := db.LedgerHeaderFirstRow()
	lastLedger := db.LedgerHeaderLastRow()

	if config.Start.Explicit {
		if config.Start.Value < 0 {
			first = lastLedger.LedgerSeq + config.Start.Value + 1
		} else if config.Start.Value > 0 {
			first = firstLedger.LedgerSeq + config.Start.Value
		}
	} else if config.Start.Value != 0 {
		first = config.Start.Value
	} else {
		first = firstLedger.LedgerSeq
	}

	if *config.Count == 0 {
		last = lastLedger.LedgerSeq
	} else {
		last = first + *config.Count - 1
	}

	return first, last
}

func createBar(count int) {
	bar = progressbar.NewOptions(
		count,
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionThrottle(500*time.Millisecond),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetWidth(100),
	)

	bar.RenderBlank()
}

func finishBar() {
	if !*config.Verbose {
		bar.Finish()
	}
}

func blockCount(count int) (blocks int) {
	blocks = count / db.LedgerHeaderRowBatchSize

	if count%db.LedgerHeaderRowBatchSize > 0 {
		blocks = blocks + 1
	}

	return blocks
}
