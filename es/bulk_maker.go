package es

import (
	"bytes"

	"github.com/astroband/astrologer/db"
)

// BulkMaker creates es bulk from ledger data
type BulkMaker struct {
	ledgers      db.LedgerHeaderRow
	transactions []db.TxHistoryRow
	fees         []db.TxFeeHistoryRow
	buf          *bytes.Buffer
}

// NewBulkMaker returns new BulkMaker structure
func NewBulkMaker(l db.LedgerHeaderRow, t []db.TxHistoryRow, f []db.TxFeeHistoryRow, b *bytes.Buffer) *BulkMaker {
	return &BulkMaker{l, t, f, b}
}

func (m BulkMaker) Make() {

}
