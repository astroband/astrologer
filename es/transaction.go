package es

import (
	"fmt"
	"time"

	"github.com/gzigzigzeo/stellar-core-export/db"
	"github.com/stellar/go/xdr"
)

type Memo struct {
	Type  int    `json:"type"`
	Value string `json:"value"`
}

// Transaction represents ES-serializable transaction
type Transaction struct {
	ID              string    `json:"id"`
	Index           byte      `json:"idx"`
	Seq             int       `json:"seq"`
	Order           string    `json:"order"`
	Fee             int       `json:"fee"`
	FeeCharged      int       `json:"fee_charged"`
	OperationCount  int       `json:"operation_count"`
	CloseTime       time.Time `json:"close_time"`
	Successful      bool      `json:"succesful"`
	ResultCode      int       `json:"result_code"`
	SourceAccountID string    `json:"source_account_id"`

	*Memo `json:"memo,omitempty"`
}

// NewTransaction creates LedgerHeader from LedgerHeaderRow
func NewTransaction(row *db.TxHistoryRow, t time.Time) *Transaction {
	resultCode := row.Result.Result.Result.Code

	tx := &Transaction{
		ID:              row.ID,
		Index:           byte(row.Index),
		Seq:             row.LedgerSeq,
		Order:           fmt.Sprintf("%d:%d", row.LedgerSeq, row.Index),
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

	return tx
}

func (tx *Transaction) DocID() string {
	return tx.ID
}

func (tx *Transaction) IndexName() string {
	return txIndexName
}
