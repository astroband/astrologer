package commands

import (
	"fmt"
	"os"
	"strconv"

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
	table.SetHeader([]string{"Min", "Max", "Count"})

	g = append(g, first.LedgerSeq)

	for i := 0; i < len(gaps); i++ {
		gap := gaps[i]
		g = append(g, gap.Start-1)
		g = append(g, gap.End+1)
	}

	g = append(g, last.LedgerSeq)

	total := 0

	for i := 0; i < len(g)/2; i++ {
		min := g[i*2]
		max := g[i*2+1]
		count := max - min + 1

		total += count

		table.Append([]string{
			strconv.Itoa(min),
			strconv.Itoa(max),
			strconv.Itoa(count),
		})
	}

	table.SetFooter([]string{"", "Total", strconv.Itoa(total)})

	table.Render()
}
