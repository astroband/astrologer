package es

import (
	"encoding/hex"
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
)

// LedgerHeader represents json-serializable struct for LedgerHeader to index
type LedgerHeader struct {
	Hash           string      `json:"hash"`
	PrevHash       string      `json:"prev_hash"`
	BucketListHash string      `json:"bucket_list_hash"`
	Seq            int         `json:"seq"`
	PagingToken    PagingToken `json:"paging_token"`
	CloseTime      time.Time   `json:"close_time"`
	Version        int         `json:"version"`
	TotalCoins     int         `json:"total_coins"`
	FeePool        int         `json:"fee_pool"`
	InflationSeq   int         `json:"inflation_seq"`
	IDPool         int         `json:"id_pool"`
	BaseFee        int         `json:"base_fee"`
	BaseReserve    int         `json:"base_reserve"`
	MaxTxSetSize   int         `json:"max_tx_size"`
}

// NewLedgerHeader creates LedgerHeader from LedgerHeaderRow
func NewLedgerHeader(row *db.LedgerHeaderRow) *LedgerHeader {
	return &LedgerHeader{
		Hash:           row.Hash,
		PrevHash:       row.PrevHash,
		BucketListHash: row.BucketListHash,
		Seq:            row.LedgerSeq,
		PagingToken:    PagingToken{LedgerSeq: row.LedgerSeq},
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

func NewLedgerHeaderFromHistory(historyEntry xdr.LedgerHeaderHistoryEntry) *LedgerHeader {
	header := historyEntry.Header
	seq := int(header.LedgerSeq)
	pagingToken := PagingToken{LedgerSeq: seq}

	return &LedgerHeader{
		ID:             pagingToken.String(),
		Hash:           hex.EncodeToString(historyEntry.Hash[:]),
		PrevHash:       hex.EncodeToString(header.PreviousLedgerHash[:]),
		BucketListHash: hex.EncodeToString(header.BucketListHash[:]),
		Seq:            seq,
		PagingToken:    pagingToken,
		CloseTime:      time.Unix(int64(header.ScpValue.CloseTime), 0),
		Version:        int(header.LedgerVersion),
		TotalCoins:     int(header.TotalCoins),
		FeePool:        int(header.FeePool),
		InflationSeq:   int(header.InflationSeq),
		IDPool:         int(header.IdPool),
		BaseFee:        int(header.BaseFee),
		BaseReserve:    int(header.BaseReserve),
		MaxTxSetSize:   int(header.MaxTxSetSize),
	}
}

// DocID returns es id (seq number in this case)
func (h *LedgerHeader) DocID() *string {
	s := h.PagingToken.String()
	return &s
}

// IndexName returns index name
func (h *LedgerHeader) IndexName() IndexName {
	return ledgerHeaderIndexName
}
