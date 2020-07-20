package commands

import (
	"bytes"
	log "github.com/sirupsen/logrus"

	"github.com/astroband/astrologer/es"
	lb "github.com/stellar/go/exp/ingest/ledgerbackend"
)

// ExportCommandConfig represents configuration options for `export` CLI command
type ExportCommandConfig struct {
	Start             int
	Count             int
	RetryCount        int
	DryRun            bool
	BatchSize         int
	NetworkPassphrase string
}

// ExportCommand represents the `export` CLI command
type ExportCommand struct {
	ES     es.Adapter
	Config ExportCommandConfig

	firstLedger uint32
	lastLedger  uint32
}

// Execute starts the export process
func (cmd *ExportCommand) Execute() {
	total := cmd.Config.Count

	cmd.firstLedger = uint32(cmd.Config.Start)
	cmd.lastLedger = cmd.firstLedger + uint32(cmd.Config.Count) - 1

	if total == 0 {
		log.Fatal("Nothing to export within given range!", cmd.firstLedger, cmd.lastLedger)
	}

	log.Infof("Exporting ledgers from %d to %d. Total: %d ledgers\n", cmd.firstLedger, cmd.lastLedger, total)
	log.Infof("Will insert %d batches %d ledgers each\n", cmd.blockCount(total), cmd.Config.BatchSize)

	ledgerBackend := lb.NewCaptive(
		"stellar-core",
		cmd.Config.NetworkPassphrase,
		getHistoryURLs(cmd.Config.NetworkPassphrase),
	)

	err := ledgerBackend.PrepareRange(cmd.firstLedger, cmd.lastLedger)

	if err != nil {
		log.Fatal(err)
	}

	var batchBuffer bytes.Buffer

	for ledgerSeq := cmd.firstLedger; ledgerSeq <= cmd.lastLedger; ledgerSeq++ {
		_, meta, err := ledgerBackend.GetLedger(ledgerSeq)

		if err != nil {
			// FIXME skip instead of failing
			log.Fatal(err)
		}

		es.SerializeLedgerFromHistory(meta, &batchBuffer)

		if (ledgerSeq-cmd.firstLedger+1)%uint32(cmd.Config.BatchSize) == 0 || ledgerSeq == cmd.lastLedger {
			payload := batchBuffer.String()
			pool.Submit(func() {
				log.Printf("Gonna bulk insert %d bytes\n", len(payload))
				err := cmd.ES.BulkInsert(payload)

				if err != nil {
					log.Fatal("Cannot bulk insert", err)
				} else {
					log.Printf("Batch successfully inserted\n")
				}
			})

			batchBuffer.Reset()
		}
	}

	// for i := 0; i < cmd.blockCount(total); i++ {
	// 	var b bytes.Buffer
	// 	ledgerCounter := 0
	// 	batchNum := i + 1

	// 	for meta := range channel {
	// 		seq := int(meta.V0.LedgerHeader.Header.LedgerSeq)

	// 		if seq < cmd.firstLedger || seq > cmd.lastLedger {
	// 			continue
	// 		}

	// 		ledgerCounter += 1

	// 		log.Println(seq)

	// 		es.SerializeLedgerFromHistory(meta, &b)

	// 		log.Printf("Ledger %d of %d in batch %d\n", ledgerCounter, cmd.Config.BatchSize, batchNum)

	// 		if ledgerCounter == cmd.Config.BatchSize {
	// 			break
	// 		}
	// 	}

	// 	if cmd.Config.DryRun {
	// 		continue
	// 	}

	// 	pool.Submit(func() {
	// 		log.Printf("Gonna bulk insert %d bytes\n", b.Len())
	// 		err := cmd.ES.BulkInsert(b)

	// 		if err != nil {
	// 			log.Fatal("Cannot bulk insert", err)
	// 		} else {
	// 			log.Printf("Batch %d successfully inserted\n", batchNum)
	// 		}
	// 	})
	// }

	pool.StopWait()
}

func (cmd *ExportCommand) blockCount(count int) (blocks int) {
	blocks = count / cmd.Config.BatchSize

	if count%cmd.Config.BatchSize > 0 {
		blocks = blocks + 1
	}

	return blocks
}

func getHistoryURLs(networkPassphrase string) []string {
	switch networkPassphrase {
	case "Public Global Stellar Network ; September 2015":
		return []string{
			"https://history.stellar.org/prd/core-live/core_live_001",
			"https://history.stellar.org/prd/core-live/core_live_002",
			"https://history.stellar.org/prd/core-live/core_live_003",
		}
	case "Test SDF Network ; September 2015":
		return []string{
			"http://history.stellar.org/prd/core-testnet/core_testnet_001",
			"http://history.stellar.org/prd/core-testnet/core_testnet_002",
			"http://history.stellar.org/prd/core-testnet/core_testnet_003",
		}
	default:
		return []string{}
	}
}
