package commands

import (
	"fmt"
	"log"

	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/gammazero/workerpool"
)

var (
	pool = workerpool.New(*config.Concurrency)
)

// Export command
func Export() {
	first, last := getRange()
	total := db.LedgerHeaderRowCount(first, last)

	if total == 0 {
		log.Fatal("Nothing to export within given range!", first, last)
	}

	fmt.Println("Exporting ledgers from", first, "to", last, "total", total)
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

// func (*Export) Do() {
// count := db.LedgerHeaderRowCount(*config.Start, *config.Count)

// if count == 0 {
// 	log.Fatal("Nothing to export!")
// }

// bar := progressbar.NewOptions(
// 	count,
// 	progressbar.OptionEnableColorCodes(true),
// 	progressbar.OptionShowCount(),
// 	progressbar.OptionThrottle(500*time.Millisecond),
// 	progressbar.OptionSetRenderBlankState(true),
// 	progressbar.OptionSetWidth(100),
// )

// bar.RenderBlank()

// blocks := count / db.LedgerHeaderRowBatchSize
// if count%db.LedgerHeaderRowBatchSize > 0 {
// 	blocks = blocks + 1
// }

// for i := 0; i < blocks; i++ {
// 	i := i
// 	pool.Submit(func() { fetch(i, bar) })
// }

// pool.StopWait()

// if !*config.Verbose {
// 	bar.Finish()
// }

// fmt.Println("Finished!")
// }
