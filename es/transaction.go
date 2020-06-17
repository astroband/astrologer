package es

import (
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/util"
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
}

// NewTransaction creates LedgerHeader from LedgerHeaderRow
func NewTransaction(row *db.TxHistoryRow, t time.Time) (*Transaction, error) {
	resultCode := row.Result.Result.Result.Code

	var (
		err             error
		fee             int
		operationCount  int
		sourceAccountId string
		memo            xdr.Memo
		timeBounds      *xdr.TimeBounds
	)

	switch row.Envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		tx := row.Envelope.V0.Tx
		fee = int(tx.Fee)
		operationCount = len(tx.Operations)
		sourceAccountId, err = util.EncodeEd25519(tx.SourceAccountEd25519)
		if err != nil {
			return nil, err
		}
		memo = tx.Memo
		timeBounds = tx.TimeBounds
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		tx := row.Envelope.V1.Tx
		fee = int(tx.Fee)
		operationCount = len(tx.Operations)
		sourceAccountId, err = util.EncodeMuxedAccount(tx.SourceAccount)

		if err != nil {
			return nil, err
		}
		memo = tx.Memo
		timeBounds = tx.TimeBounds
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		tx := row.Envelope.FeeBump.Tx
		innerTx := tx.InnerTx.V1.Tx
		fee = int(innerTx.Fee)
		operationCount = len(innerTx.Operations)
		sourceAccountId, err = util.EncodeMuxedAccount(innerTx.SourceAccount)

		if err != nil {
			return nil, err
		}
		memo = innerTx.Memo
		timeBounds = innerTx.TimeBounds
	}

	tx := &Transaction{
		ID:              row.ID,
		Index:           row.Index,
		Seq:             row.LedgerSeq,
		PagingToken:     PagingToken{LedgerSeq: row.LedgerSeq, TransactionOrder: row.Index},
		Fee:             fee,
		FeeCharged:      int(row.Result.Result.FeeCharged),
		OperationCount:  operationCount,
		CloseTime:       t,
		Successful:      resultCode == xdr.TransactionResultCodeTxSuccess,
		ResultCode:      int(resultCode),
		SourceAccountID: sourceAccountId,
	}

	if memo.Type != xdr.MemoTypeMemoNone {
		value := row.MemoValue()

		tx.Memo = &Memo{
			Type:  int(memo.Type),
			Value: value.String,
		}
	}

	if timeBounds != nil {
		tx.TimeBounds = &TimeBounds{
			MinTime: int64(timeBounds.MinTime),
			MaxTime: int64(timeBounds.MaxTime),
		}
	}

	return tx, nil
}

// DocID return es transaction id (tx id in this case)
func (tx *Transaction) DocID() *string {
	return &tx.ID
}

// IndexName returns tx index name
func (tx *Transaction) IndexName() IndexName {
	return txIndexName
}
