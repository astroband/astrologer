package es

import (
	"bytes"

	"github.com/astroband/astrologer/db"
)

// BulkMaker creates es bulk from ledger data
type BulkMaker struct {
	ledgerRow       db.LedgerHeaderRow
	ledgerHeader    *LedgerHeader
	transactionRows []db.TxHistoryRow
	transactions    []*Transaction
	fees            []db.TxFeeHistoryRow
	buffer          *bytes.Buffer
}

// NewBulkMaker returns new BulkMaker structure
func NewBulkMaker(l db.LedgerHeaderRow, t []db.TxHistoryRow, f []db.TxFeeHistoryRow, b *bytes.Buffer) *BulkMaker {
	h := NewLedgerHeader(&l)

	txs := make([]*Transaction, len(t))
	for i := 0; i < len(t); i++ {
		txs[i] = NewTransaction(&t[i], h.CloseTime)
	}

	return &BulkMaker{
		ledgerRow:       l,
		ledgerHeader:    h,
		transactionRows: t,
		transactions:    txs,
		fees:            f,
		buffer:          b,
	}
}

// Make creates bulk
func (m *BulkMaker) Make() {
	m.makeLedger()
	m.makeTransactions()
}

func (m *BulkMaker) makeLedger() {
	SerializeForBulk(m.ledgerHeader, m.buffer)
}

func (m *BulkMaker) makeTransactions() {
	for _, transaction := range m.transactions {
		SerializeForBulk(transaction, m.buffer)
	}
}

// for t := 0; t < len(txs); t++ {
// 	var metas []xdr.OperationMeta

// 	txRow := &txs[t]
// 	ops := txRow.Envelope.Tx.Operations
// 	results := txRow.Result.Result.Result.Results

// 	if v1, ok := txRow.Meta.GetV1(); ok {
// 		metas = v1.Operations
// 	} else {
// 		metas, ok = txRow.Meta.GetOperations()
// 	}

// 	tx := NewTransaction(txRow, h.CloseTime)
// 	SerializeForBulk(tx, b)

// 	for o := 0; o < len(ops); o++ {
// 		op := NewOperation(tx, &ops[o], byte(o))

// 		if results != nil {
// 			AppendResult(op, &(*results)[o])
// 		}

// 		SerializeForBulk(op, b)
// 	}

// 	for o := 0; o < len(metas); o++ {
// 		pagingToken := PagingToken{LedgerSeq: h.Seq, TransactionOrder: tx.Index, OperationOrder: uint8(o + 1), AuxOrder1: 1}
// 		bl := NewBalanceExtractor(metas[o].Changes, h.CloseTime, BalanceSourceMeta, pagingToken).Extract()

// 		for _, balance := range bl {
// 			SerializeForBulk(balance, b)
// 		}
// 	}
// }
