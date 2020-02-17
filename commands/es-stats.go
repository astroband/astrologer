package commands

import (
	"os"
	"strconv"

	"github.com/astroband/astrologer/es"
	"github.com/olekukonko/tablewriter"
)

const step = 10000

type EsStatsCommand struct {
	ES es.Adapter
}

// EsStats prints ledger statistics for current database
func (cmd *EsStatsCommand) Execute() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"From", "To", "Doc_count"})

	min, max := cmd.ES.MinMaxSeq()
	buckets := cmd.esRanges(min, max)

	for i := 0; i < len(buckets); i++ {
		bucket := buckets[i].(map[string]interface{})
		count := int(bucket["doc_count"].(float64))
		from := int(bucket["from"].(float64))
		to := int(bucket["to"].(float64))

		table.Append([]string{
			strconv.Itoa(from),
			strconv.Itoa(to),
			strconv.Itoa(count),
		})
	}

	table.Render()
}

func (cmd *EsStatsCommand) esRanges(min int, max int) []interface{} {
	var ranges []map[string]interface{}

	for i := min; i < max; i += step {
		to := i + step
		if to > max {
			to = max
		}
		ranges = append(ranges, map[string]interface{}{"from": i, "to": to})
	}

	aggs := cmd.ES.LedgerSeqRangeQuery(ranges)
	buckets := aggs["buckets"].([]interface{})

	return buckets
}
