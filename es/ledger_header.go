package es

import (
	"strconv"
	"time"

	"github.com/gzigzigzeo/stellar-core-export/db"
)

// LedgerHeader represents json-serializable struct for LedgerHeader to index
type LedgerHeader struct {
	Hash           string    `json:"hash"`
	PrevHash       string    `json:"prev_hash"`
	BucketListHash string    `json:"bucket_list_hash"`
	Seq            int       `json:"seq"`
	CloseTime      time.Time `json:"close_time"`
	Version        int       `json:"version"`
	TotalCoins     int       `json:"total_coins"`
	FeePool        int       `json:"fee_pool"`
	InflationSeq   int       `json:"inflation_seq"`
	IDPool         int       `json:"id_pool"`
	BaseFee        int       `json:"base_fee"`
	BaseReserve    int       `json:"base_reserve"`
	MaxTxSetSize   int       `json:"max_tx_size"`
}

// NewLedgerHeader creates LedgerHeader from LedgerHeaderRow
func NewLedgerHeader(row *db.LedgerHeaderRow) *LedgerHeader {
	return &LedgerHeader{
		Hash:           row.Hash,
		PrevHash:       row.PrevHash,
		BucketListHash: row.BucketListHash,
		Seq:            row.LedgerSeq,
		CloseTime:      time.Unix(row.CloseTime, 0),
		Version:        int(row.Data.LedgerVersion),
		TotalCoins:     int(row.Data.TotalCoins),
		FeePool:        int(row.Data.FeePool),
		InflationSeq:   int(row.Data.InflationSeq),
		IDPool:         int(row.Data.IdPool),
		BaseFee:        int(row.Data.BaseFee),
		BaseReserve:    int(row.Data.BaseReserve),
		MaxTxSetSize:   int(row.Data.MaxTxSetSize),
	}
}

func (h *LedgerHeader) DocID() string {
	return strconv.Itoa(h.Seq)
}

func (h *LedgerHeader) IndexName() string {
	return ledgerHeaderIndexName
}
