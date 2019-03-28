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
	Index           int       `json:"idx"`
	Seq             int       `json:"seq"`
	Order           string    `json:"order"`
	Fee             int       `json:"fee"`
	FeePaid         int       `json:"fee_paid"`
	OperationCount  int       `json:"operation_count"`
	CloseTime       time.Time `json:"close_time"`
	Successful      bool      `json:"succesful"`
	ResultCode      int       `json:"result_code"`
	SourceAccountID string    `json:"source_account"`

	*Memo `json:"memo,omitempty"`
}

// NewTransaction creates LedgerHeader from LedgerHeaderRow
func NewTransaction(row *db.TxHistoryRow, t time.Time) *Transaction {
	var e xdr.TransactionEnvelope
	var r xdr.TransactionResult

	xdr.SafeUnmarshalBase64(row.TxBody, &e)
	xdr.SafeUnmarshalBase64(row.TxResult, &r)

	tx := &Transaction{
		ID:              row.TxID,
		Index:           row.TxIndex,
		Seq:             row.LedgerSeq,
		Order:           fmt.Sprintf("%d:%d", row.LedgerSeq, row.TxIndex),
		Fee:             int(e.Tx.Fee),
		FeePaid:         int(r.FeeCharged),
		OperationCount:  len(e.Tx.Operations),
		CloseTime:       t,
		Successful:      r.Result.Code == xdr.TransactionResultCodeTxSuccess,
		ResultCode:      int(r.Result.Code),
		SourceAccountID: e.Tx.SourceAccount.Address(),
	}

	tx.setMemo(&e)

	return tx
}

func (tx *Transaction) setMemo(e *xdr.TransactionEnvelope) {
	if e.Tx.Memo.Type != xdr.MemoTypeMemoNone {
		var v string = ""

		// switch e.Tx.Memo.Type {
		// case xdr.MemoTypeMemoHash:
		// 	// b, err := e.Tx.Memo.Hash.MarshalBinary()
		// 	// if err != nil {
		// 	// 	log.Fatal(err)
		// 	// }
		// 	v = hex.EncodeToString(([]byte)e.Tx.Memo.Hash)
		// case xdr.MemoTypeMemoReturn:

		// }

		tx.Memo = &Memo{
			Type:  int(e.Tx.Memo.Type),
			Value: v,
		}
	}
}

func (tx *Transaction) DocID() string {
	return tx.ID
}

func (tx *Transaction) IndexName() string {
	return txIndexName
}
