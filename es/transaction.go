package es

import (
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
)

type Memo struct {
	Type  int    `json:"type"`
	Value string `json:"value"`
}

type TimeBounds struct {
	MinTime int64 `json:"min_time"`
	MaxTime int64 `json:"max_time"`
}

// Transaction represents ES-serializable transaction
type Transaction struct {
	ID              string    `json:"id"`
	Index           byte      `json:"idx"`
	Seq             int       `json:"seq"`
	Order           int       `json:"order"`
	Fee             int       `json:"fee"`
	FeeCharged      int       `json:"fee_charged"`
	OperationCount  int       `json:"operation_count"`
	CloseTime       time.Time `json:"close_time"`
	Successful      bool      `json:"succesful"`
	ResultCode      int       `json:"result_code"`
	SourceAccountID string    `json:"source_account_id"`

	*TimeBounds `json:"time_bounds,omitempty"`
	*Memo       `json:"memo,omitempty"`
}

// NewTransaction creates LedgerHeader from LedgerHeaderRow
func NewTransaction(row *db.TxHistoryRow, t time.Time) *Transaction {
	resultCode := row.Result.Result.Result.Code

	tx := &Transaction{
		ID:              row.ID,
		Index:           byte(row.Index),
		Seq:             row.LedgerSeq,
		Order:           row.LedgerSeq*100 + row.Index,
		Fee:             int(row.Envelope.Tx.Fee),
		FeeCharged:      int(row.Result.Result.FeeCharged),
		OperationCount:  len(row.Envelope.Tx.Operations),
		CloseTime:       t,
		Successful:      resultCode == xdr.TransactionResultCodeTxSuccess,
		ResultCode:      int(resultCode),
		SourceAccountID: row.Envelope.Tx.SourceAccount.Address(),
	}

	if row.Envelope.Tx.Memo.Type != xdr.MemoTypeMemoNone {
		value := row.MemoValue()

		tx.Memo = &Memo{
			Type:  int(row.Envelope.Tx.Memo.Type),
			Value: value.String,
		}
	}

	if row.Envelope.Tx.TimeBounds != nil {
		tx.TimeBounds = &TimeBounds{
			MinTime: int64(row.Envelope.Tx.TimeBounds.MinTime),
			MaxTime: int64(row.Envelope.Tx.TimeBounds.MaxTime),
		}
	}

	return tx
}

// DocID return es transaction id (tx id in this case)
func (tx *Transaction) DocID() *string {
	return &tx.ID
}

// IndexName returns tx index name
func (tx *Transaction) IndexName() string {
	return txIndexName
}
