package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/astroband/astrologer/config"
	"github.com/astroband/astrologer/db"
	"github.com/olekukonko/tablewriter"
)

// Stats prints ledger statistics for current database
func Stats() {
	var g []int

	first := db.LedgerHeaderFirstRow()
	last := db.LedgerHeaderLastRow()
	gaps := db.LedgerHeaderGaps()

	if (first == nil) || (last == nil) {
		fmt.Println("Current database is empty!")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Min", "Max", "Count", "ES"})

	g = append(g, first.LedgerSeq)

	for i := 0; i < len(gaps); i++ {
		gap := gaps[i]
		g = append(g, gap.Start-1)
		g = append(g, gap.End+1)
	}

	g = append(g, last.LedgerSeq)

	total := 0
	totalES := 0

	for i := 0; i < len(g)/2; i++ {
		min := g[i*2]
		max := g[i*2+1]
		count := max - min + 1
		countES := esCount(min, max)

		total += count
		totalES += countES

		table.Append([]string{
			strconv.Itoa(min),
			strconv.Itoa(max),
			strconv.Itoa(count),
			strconv.Itoa(countES),
		})
	}

	table.SetFooter([]string{"", "Total", strconv.Itoa(total), strconv.Itoa(totalES)})

	table.Render()
}

func esCount(min int, max int) int {
	var r map[string]interface{}
	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"seq": map[string]interface{}{
					"gte": min,
					"lte": max,
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := config.ES.Count(
		config.ES.Count.WithIndex("ledger"),
		config.ES.Count.WithBody(&buf),
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

	return int(r["count"].(float64))
}
