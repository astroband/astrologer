package commands

import (
	"github.com/astroband/astrologer/es"
  "log"
)

type DedupeCommandConfig struct {
	DryRun bool
  Start int
  Count int
}

type DedupeCommand struct {
	ES     es.Adapter
	Config DedupeCommandConfig
}

func (cmd *DedupeCommand) Execute() {
  indexName := "ledger"

  log.Printf("Searching '%s' index...\n", indexName)

  var duplicates []string
  var after *string

  removeCount := 0

  for {
    buckets := cmd.ES.GroupLedgersBySeq(
      cmd.Config.Start,
      cmd.Config.Start + cmd.Config.Count,
      after,
    )

    if len(buckets) == 0 {
      break
    }

    for _, bucket := range buckets {
      if bucket.DocCount == 1 {
        continue
      }

      removeCount += bucket.DocCount - 1
      duplicates = append(duplicates, bucket.FieldValue)
    }

    after = &buckets[len(buckets) - 1].FieldValue
    // log.Println(*after)
  }

  log.Println(len(duplicates), removeCount)
}
