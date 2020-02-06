package commands

import (
	// "log"
	"os"
	"strconv"

	"github.com/astroband/astrologer/es"
	"github.com/olekukonko/tablewriter"
)

const step = 10000

// EsStats prints ledger statistics for current database
func EsStats() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"From", "To", "Doc_count"})

	min, max := es.Adapter.MinMaxSeq()
	buckets := esRanges(min, max)

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

func esRanges(min int, max int) []interface{} {
	var ranges []map[string]interface{}

	for i := min; i < max; i += step {
		to := i + step
		if to > max {
			to = max
		}
		ranges = append(ranges, map[string]interface{}{"from": i, "to": to})
	}

	aggs := es.Adapter.LedgerSeqRangeQuery(ranges)
	buckets := aggs["buckets"].([]interface{})

	return buckets
}
