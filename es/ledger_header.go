package es

import (
	"strconv"
	"time"

	"github.com/gzigzigzeo/stellar-core-export/db"
	"github.com/stellar/go/xdr"
)

// LedgerHeader represents json-serializable struct for LedgerHeader to index
type LedgerHeader struct {
	DocID          string    `json:"-"`
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
	var h xdr.LedgerHeader

	xdr.SafeUnmarshalBase64(row.Data, &h)

	return &LedgerHeader{
		DocID:          strconv.Itoa(row.LedgerSeq),
		Hash:           row.Hash,
		PrevHash:       row.PrevHash,
		BucketListHash: row.BucketListHash,
		Seq:            row.LedgerSeq,
		CloseTime:      time.Unix(row.CloseTime, 0),
		Version:        int(h.LedgerVersion),
		TotalCoins:     int(h.TotalCoins),
		FeePool:        int(h.FeePool),
		InflationSeq:   int(h.InflationSeq),
		IDPool:         int(h.IdPool),
		BaseFee:        int(h.BaseFee),
		BaseReserve:    int(h.BaseReserve),
		MaxTxSetSize:   int(h.MaxTxSetSize),
	}
}
