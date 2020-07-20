package es

import (
	"encoding/hex"
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/support"
	"github.com/astroband/astrologer/util"
	"github.com/stellar/go/network"
	"github.com/stellar/go/xdr"
)

// transactionData represents all necessary pieces we need
// from stellar-core to ingest single transaction
type transactionData struct {
	ledgerSeq int
	// ledger close time
	closeTime time.Time
	// index of transaction in ledger
	index  int
	v      txVersion
	xdrV0  *xdr.TransactionV0
	xdrV1  *xdr.Transaction
	result xdr.TransactionResultPair
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

// NewTransaction creates LedgerHeader from LedgerHeaderRow
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

func NewTransactionFromXDR(data *transactionData) (*Transaction, error) {
	resultCode := data.result.Result.Result.Code

	var (
		err             error
		txHash          [32]byte
		sourceAccountID string
		maxFee          int
		operationsCount int
		memo            Memo
		timeBounds      TimeBounds
	)

	switch data.v {
	case v0:
		//FIXME remove hardcoded network passphrase
		txHash, err = network.HashTransactionV0(
			*data.xdrV0,
			"Public Global Stellar Network ; September 2015",
		)
		if err != nil {
			return nil, err
		}

		sourceAccountID, err = util.EncodeEd25519(data.xdrV0.SourceAccountEd25519)
		maxFee = int(data.xdrV0.Fee)
		operationsCount = len(data.xdrV0.Operations)

		if data.xdrV0.Memo.Type != xdr.MemoTypeMemoNone {
			value := support.MemoValue(data.xdrV0.Memo)

			memo = Memo{
				Type:  int(data.xdrV0.Memo.Type),
				Value: value.String,
			}
		}

		if data.xdrV0.TimeBounds != nil {
			timeBounds = TimeBounds{
				MinTime: int64(data.xdrV0.TimeBounds.MinTime),
				MaxTime: int64(data.xdrV0.TimeBounds.MaxTime),
			}
		}

	case v1:
		//FIXME remove hardcoded network passphrase
		txHash, err = network.HashTransaction(
			*data.xdrV1,
			"Public Global Stellar Network ; September 2015",
		)
		if err != nil {
			return nil, err
		}

		sourceAccountID, err = util.EncodeMuxedAccount(data.xdrV1.SourceAccount)
		maxFee = int(data.xdrV1.Fee)
		operationsCount = len(data.xdrV1.Operations)

		if data.xdrV1.Memo.Type != xdr.MemoTypeMemoNone {
			value := support.MemoValue(data.xdrV1.Memo)

			memo = Memo{
				Type:  int(data.xdrV1.Memo.Type),
				Value: value.String,
			}
		}

		if data.xdrV1.TimeBounds != nil {
			timeBounds = TimeBounds{
				MinTime: int64(data.xdrV1.TimeBounds.MinTime),
				MaxTime: int64(data.xdrV1.TimeBounds.MaxTime),
			}
		}
	}

	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		ID:              hex.EncodeToString(txHash[:]),
		Index:           data.index,
		Seq:             data.ledgerSeq,
		PagingToken:     PagingToken{LedgerSeq: data.ledgerSeq, TransactionOrder: data.index},
		MaxFee:          maxFee,
		FeeCharged:      int(data.result.Result.FeeCharged),
		OperationCount:  operationsCount,
		CloseTime:       data.closeTime,
		Successful:      resultCode == xdr.TransactionResultCodeTxSuccess,
		ResultCode:      int(resultCode),
		SourceAccountID: sourceAccountID,
		Memo:            &memo,
		TimeBounds:      &timeBounds,
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
