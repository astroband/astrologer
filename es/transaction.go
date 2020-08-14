package es

import (
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/support"
	"github.com/stellar/go/xdr"
)

// transactionData represents all necessary pieces we need
// from stellar-core to ingest single transaction
type transactionData struct {
	id        string
	ledgerSeq int
	// ledger close time
	closeTime time.Time
	// index of transaction in ledger
	index      int
	xdrV0      *xdr.TransactionV0
	xdrV1      *xdr.Transaction
	feeBump    *xdr.FeeBumpTransaction
	result     xdr.TransactionResultPair
	feeAccount string
	maxFee     int
}

// Transaction represents ES-serializable transaction
type Transaction struct {
	ID              string      `json:"id"`
	Index           int         `json:"idx"`
	Seq             int         `json:"seq"`
	PagingToken     PagingToken `json:"paging_token"`
	MaxFee          int         `json:"max_fee"`
	FeeCharged      int         `json:"fee_charged"`
	FeeAccountID    string      `json:"fee_account_id"`
	OperationCount  int         `json:"operation_count"`
	CloseTime       time.Time   `json:"close_time"`
	Successful      bool        `json:"successful"`
	ResultCode      int         `json:"result_code"`
	SourceAccountID string      `json:"source_account_id"`

	*TimeBounds `json:"time_bounds,omitempty"`
	*Memo       `json:"memo,omitempty"`
}

// NewTransaction creates Transaction from TxHistoryRow
func (s *ledgerSerializer) NewTransaction(row *db.TxHistoryRow, t time.Time) (*Transaction, error) {
	var (
		err      error
		envelope = row.Envelope
		result   = row.Result.Result.Result
		success  bool
	)

	if envelope.IsFeeBump() {
		success = result.Code == xdr.TransactionResultCodeTxFeeBumpInnerSuccess
	} else {
		success = result.Code == xdr.TransactionResultCodeTxSuccess
	}

	accountId := envelope.SourceAccount().ToAccountId()
	sourceAccountAddress, err := (&accountId).GetAddress()

	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		ID:              row.ID,
		Index:           row.Index,
		Seq:             row.LedgerSeq,
		MaxFee:          int(envelope.Fee()),
		PagingToken:     PagingToken{LedgerSeq: row.LedgerSeq, TransactionOrder: row.Index},
		FeeCharged:      int(row.Result.Result.FeeCharged),
		CloseTime:       t,
		Successful:      success,
		ResultCode:      int(result.Code),
		SourceAccountID: sourceAccountAddress,
		OperationCount:  len(envelope.Operations()),
	}

	if envelope.IsFeeBump() {
		feeSourceAccountId := envelope.FeeBumpAccount().ToAccountId()
		feeSourceAddress, err := (&feeSourceAccountId).GetAddress()

		if err != nil {
			return nil, err
		}

		transaction.FeeAccountID = feeSourceAddress
		transaction.MaxFee = int(envelope.FeeBumpFee())
	}

	if envelope.Memo().Type != xdr.MemoTypeMemoNone {
		transaction.Memo = &Memo{
			Type:  int(envelope.Memo().Type),
			Value: support.MemoValue(row.Envelope.Memo()).String,
		}
	}

	if envelope.TimeBounds() != nil {
		transaction.TimeBounds = &TimeBounds{
			MinTime: int64(envelope.TimeBounds().MinTime),
			MaxTime: int64(envelope.TimeBounds().MaxTime),
		}
	}

	return transaction, nil
}

// DocID return es transaction id (tx id in this case)
func (tx *Transaction) DocID() *string {
	return &tx.ID
}

// IndexName returns tx index name
func (tx *Transaction) IndexName() IndexName {
	return txIndexName
}
