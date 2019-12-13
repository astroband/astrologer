package es

import (
	"encoding/base64"
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
)

// Transaction represents ES-serializable transaction
type Transaction struct {
	ID              string      `json:"id"`
	Index           int         `json:"idx"`
	Seq             int         `json:"seq"`
	PagingToken     PagingToken `json:"paging_token"`
	Fee             int         `json:"fee"`
	FeeCharged      int         `json:"fee_charged"`
	OperationCount  int         `json:"operation_count"`
	CloseTime       time.Time   `json:"close_time"`
	Successful      bool        `json:"successful"`
	ResultCode      int         `json:"result_code"`
	SourceAccountID string      `json:"source_account_id"`

	*TimeBounds `json:"time_bounds,omitempty"`
	*Memo       `json:"memo,omitempty"`

	Meta    string `json:"meta"`
	FeeMeta string `json:"fee_meta"`
}

// NewTransaction creates Transaction from TxHistoryRow
func NewTransaction(row *db.TxHistoryRow, feeRow *db.TxFeeHistoryRow, t time.Time) *Transaction {
	resultCode := row.Result.Result.Result.Code
	metaBinary, err := row.Meta.MarshalBinary()
	feeMetaBinary, err := feeRow.Changes.MarshalBinary()

	if err != nil {
		//FIXME What should we do here?
	}

	tx := &Transaction{
		ID:              row.ID,
		Index:           row.Index,
		Seq:             row.LedgerSeq,
		PagingToken:     PagingToken{LedgerSeq: row.LedgerSeq, TransactionOrder: row.Index},
		Fee:             int(row.Envelope.Tx.Fee),
		FeeCharged:      int(row.Result.Result.FeeCharged),
		OperationCount:  len(row.Envelope.Tx.Operations),
		CloseTime:       t,
		Successful:      resultCode == xdr.TransactionResultCodeTxSuccess,
		ResultCode:      int(resultCode),
		SourceAccountID: row.Envelope.Tx.SourceAccount.Address(),
		Meta:            base64.StdEncoding.EncodeToString(metaBinary),
		FeeMeta:         base64.StdEncoding.EncodeToString(feeMetaBinary),
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
