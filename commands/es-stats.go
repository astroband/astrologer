package commands

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/astroband/astrologer/config"
	"github.com/olekukonko/tablewriter"
)

// EsStats prints ledger statistics for current database
func EsStats() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"From", "To", "Doc_count"})

	min, max := esMinMax()
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

func esMinMax() (min int, max int) {
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"seq_stats": map[string]interface{}{
				"stats": map[string]interface{}{
					"field": "seq",
				},
			},
		},
	}

	r := searchLedgers(query)

	aggs := r["aggregations"].(map[string]interface{})["seq_stats"].(map[string]interface{})

	min = int(aggs["min"].(float64))
	max = int(aggs["max"].(float64))

	return min, max
}

func esRanges(min int, max int) []interface{} {
	var ranges []map[string]interface{}

	for i := min; i < max; i += 1000000 {
		to := i + 1000000
		if to > max {
			to = max
		}
		ranges = append(ranges, map[string]interface{}{"from": i, "to": to})
	}

	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"seq_ranges": map[string]interface{}{
				"range": map[string]interface{}{
					"field":  "seq",
					"ranges": ranges,
				},
			},
		},
	}

	r := searchLedgers(query)

	aggs := r["aggregations"].(map[string]interface{})["seq_ranges"].(map[string]interface{})
	buckets := aggs["buckets"].([]interface{})

	return buckets
}

func searchLedgers(query map[string]interface{}) (r map[string]interface{}) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := config.ES.Search(
		config.ES.Search.WithIndex("ledger"),
		config.ES.Search.WithBody(&buf),
	)

	if err != nil {
		log.Fatal(err)
	}

	if res.IsError() {
		log.Fatal("Error in response", res.Body)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	res.Body.Close()

	return r
}
