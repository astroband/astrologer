package es

import (
	"encoding/hex"
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/stellar"
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
	xdr    xdr.Transaction
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
			Value: row.MemoValue().String,
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

	//FIXME remove hardcoded network passphrase
	txHash, err := network.HashTransaction(
		data.xdr,
		"Public Global Stellar Network ; September 2015",
	)

	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		ID:              hex.EncodeToString(txHash[:]),
		Index:           data.index,
		Seq:             data.ledgerSeq,
		PagingToken:     PagingToken{LedgerSeq: data.ledgerSeq, TransactionOrder: data.index},
		Fee:             int(data.xdr.Fee),
		FeeCharged:      int(data.result.Result.FeeCharged),
		OperationCount:  len(data.xdr.Operations),
		CloseTime:       data.closeTime,
		Successful:      resultCode == xdr.TransactionResultCodeTxSuccess,
		ResultCode:      int(resultCode),
		SourceAccountID: data.xdr.SourceAccount.Address(),
	}

	if data.xdr.Memo.Type != xdr.MemoTypeMemoNone {
		value := stellar.MemoValue(data.xdr.Memo)

		tx.Memo = &Memo{
			Type:  int(data.xdr.Memo.Type),
			Value: value.String,
		}
	}

	if data.xdr.TimeBounds != nil {
		tx.TimeBounds = &TimeBounds{
			MinTime: int64(data.xdr.TimeBounds.MinTime),
			MaxTime: int64(data.xdr.TimeBounds.MaxTime),
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
