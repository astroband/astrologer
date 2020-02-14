package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/es"
	"github.com/olekukonko/tablewriter"
)

type StatsCommand struct {
	ES es.EsAdapter
	DB db.DbAdapter
}

// Stats prints ledger statistics for current database
func (cmd *StatsCommand) Execute() {
	var g []int

	first := cmd.DB.LedgerHeaderFirstRow()
	last := cmd.DB.LedgerHeaderLastRow()
	gaps := cmd.DB.LedgerHeaderGaps()

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
		countES := cmd.ES.LedgerCountInRange(min, max)

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
