package ingest

import (
	"strconv"
	"time"

	"github.com/stellar/go/xdr"
)

var (
	ledgerShift = 32
)

// // LedgerHeader represents json-serializable struct for LedgerHeader to index
// type LedgerHeader struct {
// 	Hash           string    `json:"hash"`
// 	PrevHash       string    `json:"prev_hash"`
// 	BucketListHash string    `json:"bucket_list_hash"`
// 	Seq            int       `json:"seq"`
// 	CloseTime      time.Time `json:"close_time"`
// 	Version        int       `json:"version"`
// 	TotalCoins     int       `json:"total_coins"`
// 	FeePool        int       `json:"fee_pool"`
// 	InflationSeq   int       `json:"inflation_seq"`
// 	IDPool         int       `json:"id_pool"`
// 	BaseFee        int       `json:"base_fee"`
// 	BaseReserve    int       `json:"base_reserve"`
// 	MaxTxSetSize   int       `json:"max_tx_size"`
// }

// // NewLedgerHeader creates LedgerHeader from LedgerHeaderRow
// func NewLedgerHeader(row *db.LedgerHeaderRow) *LedgerHeader {
// 	return &LedgerHeader{
// 		Hash:           row.Hash,
// 		PrevHash:       row.PrevHash,
// 		BucketListHash: row.BucketListHash,
// 		Seq:            row.LedgerSeq,
// 		CloseTime:      time.Unix(row.CloseTime, 0),
// 		Version:        int(row.Data.LedgerVersion),
// 		TotalCoins:     int(row.Data.TotalCoins),
// 		FeePool:        int(row.Data.FeePool),
// 		InflationSeq:   int(row.Data.InflationSeq),
// 		IDPool:         int(row.Data.IdPool),
// 		BaseFee:        int(row.Data.BaseFee),
// 		BaseReserve:    int(row.Data.BaseReserve),
// 		MaxTxSetSize:   int(row.Data.MaxTxSetSize),
// 	}
// }

// // DocID returns es id (seq number in this case)
// func (h *LedgerHeader) DocID() *string {
// 	id := strconv.Itoa(h.Seq)
// 	return &id
// }

// // IndexName returns index name
// func (h *LedgerHeader) IndexName() string {
// 	return ledgerHeaderIndexName
// }

// Ledger represents LedgerHeader with some additional data
type Ledger struct {
	ID           string
	Order        int64
	CloseTime    time.Time
	Header       *xdr.LedgerHeader
	Transactions []Transaction
}

// Transaction represents ledger transaction with all nested data
type Transaction struct {
	ID         string
	Order      int64
	ResultCode int
	Successful bool
}

// Source represents source XDR structures to parse all ledger data
type Source struct {
	CloseTime time.Time
	L         *xdr.LedgerHeader
	Tx        []*xdr.TransactionEnvelope
	R         []*xdr.TransactionResultPair
	Meta      []*xdr.TransactionMeta
	Fee       []*xdr.LedgerEntryChanges
}

// NewLedger return ledger struct filled out with all nested data
func NewLedger(s Source) *Ledger {
	return &Ledger{
		ID:        strconv.Itoa(int(s.L.LedgerSeq)),
		Order:     int64(s.L.LedgerSeq) << uint(ledgerShift),
		CloseTime: s.CloseTime,
		Header:    s.L,
	}
}
