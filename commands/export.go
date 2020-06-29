package commands

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	progressbar "github.com/schollz/progressbar/v2"

	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
)

var (
	bar *progressbar.ProgressBar
)

// ExportCommandConfig represents configuration options for `export` CLI command
type ExportCommandConfig struct {
	Start      config.NumberWithSign
	Count      int
	RetryCount int
	DryRun     bool
	BatchSize  int
}

// ExportCommand represents the `export` CLI command
type ExportCommand struct {
	ES     es.Adapter
	DB     db.Adapter
	Config ExportCommandConfig

	firstLedger int
	lastLedger  int
}

// Execute starts the export process
func (cmd *ExportCommand) Execute() {
	cmd.firstLedger, cmd.lastLedger = cmd.getRange()

	total := cmd.DB.LedgerHeaderRowCount(cmd.firstLedger, cmd.lastLedger)

	if total == 0 {
		log.Fatal("Nothing to export within given range!", cmd.firstLedger, cmd.lastLedger)
	}

	log.Println("Exporting ledgers from", cmd.firstLedger, "to", cmd.lastLedger, "total", total)

	createBar(total)

	for i := 0; i < cmd.blockCount(total); i++ {
		i := i
		pool.Submit(func() { cmd.exportBlock(i) })
	}

	pool.StopWait()
	finishBar()
}

func (cmd *ExportCommand) exportBlock(i int) {
	var b bytes.Buffer

	rows := cmd.DB.LedgerHeaderRowFetchBatch(i, cmd.firstLedger, cmd.Config.BatchSize)

	for n := 0; n < len(rows); n++ {
		txs := cmd.DB.TxHistoryRowForSeq(rows[n].LedgerSeq)
		fees := cmd.DB.TxFeeHistoryRowsForRows(txs)

		err := es.SerializeLedger(rows[n], txs, fees, &b)

		if err != nil {
			log.Fatalf("Failed to ingest ledger %d: %v\n", rows[n].LedgerSeq, err)
		}

		if !*config.Verbose {
			bar.Add(1)
		}
	}

	if *config.Verbose {
		log.Println(b.String())
	}

	if !cmd.Config.DryRun {
		cmd.ES.IndexWithRetries(&b, cmd.Config.RetryCount)
	}
}

func (cmd *ExportCommand) index(b *bytes.Buffer, retry int) {
	indexed := cmd.ES.BulkInsert(b)

	if !indexed {
		if retry > cmd.Config.RetryCount {
			log.Fatal("Retries for bulk failed, aborting")
		}

		delay := time.Duration((rand.Intn(10) + 5))
		time.Sleep(delay * time.Second)

		cmd.index(b, retry+1)
	}
}

// Parses range of export command
func (cmd *ExportCommand) getRange() (first int, last int) {
	firstLedger := cmd.DB.LedgerHeaderFirstRow()
	lastLedger := cmd.DB.LedgerHeaderLastRow()

	if cmd.Config.Start.Explicit {
		if cmd.Config.Start.Value < 0 {
			first = lastLedger.LedgerSeq + cmd.Config.Start.Value + 1
		} else if config.Start.Value > 0 {
			first = firstLedger.LedgerSeq + cmd.Config.Start.Value
		}
	} else if cmd.Config.Start.Value != 0 {
		first = cmd.Config.Start.Value
	} else {
		first = firstLedger.LedgerSeq
	}

	if cmd.Config.Count == 0 {
		last = lastLedger.LedgerSeq
	} else {
		last = first + cmd.Config.Count - 1
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

func (cmd *ExportCommand) blockCount(count int) (blocks int) {
	blocks = count / cmd.Config.BatchSize

	if count%cmd.Config.BatchSize > 0 {
		blocks = blocks + 1
	}

	return blocks
}
