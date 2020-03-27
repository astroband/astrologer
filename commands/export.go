package commands

import (
	"bytes"
	log "github.com/sirupsen/logrus"

	"github.com/astroband/astrologer/es"
	"github.com/astroband/astrologer/stellar"
)

// ExportCommandConfig represents configuration options for `export` CLI command
type ExportCommandConfig struct {
	Start      int
	Count      int
	RetryCount int
	DryRun     bool
	BatchSize  int
}

// ExportCommand represents the `export` CLI command
type ExportCommand struct {
	ES     es.Adapter
	Config ExportCommandConfig

	firstLedger int
	lastLedger  int
}

// Execute starts the export process
func (cmd *ExportCommand) Execute() {
	total := cmd.Config.Count

	cmd.firstLedger = cmd.Config.Start
	cmd.lastLedger = cmd.firstLedger + cmd.Config.Count

	if total == 0 {
		log.Fatal("Nothing to export within given range!", cmd.firstLedger, cmd.lastLedger)
	}

	log.Infof("Exporting ledgers from %d to %d. Total: %d ledgers\n", cmd.firstLedger, cmd.lastLedger, total)
	log.Infof("Will insert %d batches %d ledgers each\n", cmd.blockCount(total), cmd.Config.BatchSize)

	for i := 0; i < cmd.blockCount(total); i++ {
		var b bytes.Buffer
		ledgerCounter := 0
		batchNum := i + 1

		for meta := range stellar.StreamLedgers(cmd.firstLedger, cmd.lastLedger) {
			seq := int(meta.V0.LedgerHeader.Header.LedgerSeq)

			if seq < cmd.firstLedger || seq > cmd.lastLedger {
				continue
			}

			ledgerCounter += 1

			log.Println(seq)

			es.SerializeLedgerFromHistory(meta, &b)

			log.Printf("Ledger %d of %d in batch %d\n", ledgerCounter, cmd.Config.BatchSize, batchNum)

			if ledgerCounter == cmd.Config.BatchSize {
				break
			}
		}

		if cmd.Config.DryRun {
			continue
		}

		pool.Submit(func() {
			log.Printf("Gonna bulk insert %d bytes\n", b.Len())
			err := cmd.ES.BulkInsert(b)

			if err != nil {
				log.Fatal("Cannot bulk insert", err)
			} else {
				log.Printf("Batch %d successfully inserted\n", batchNum)
			}
		})
	}

	pool.StopWait()
}

func (cmd *ExportCommand) blockCount(count int) (blocks int) {
	blocks = count / cmd.Config.BatchSize

	if count%cmd.Config.BatchSize > 0 {
		blocks = blocks + 1
	}

	return blocks
}
