package es

import (
	"bytes"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
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
	m.makeOperationsWithResults()
	m.makeBalancesFromMetas()
}

func (m *BulkMaker) makeLedger() {
	SerializeForBulk(m.ledgerHeader, m.buffer)
}

func (m *BulkMaker) makeTransactions() {
	for _, transaction := range m.transactions {
		SerializeForBulk(transaction, m.buffer)
	}
}

func (m *BulkMaker) makeOperationsWithResults() {
	for tIndex, t := range m.transactions {
		row := m.transactionRows[tIndex]
		operations := row.Envelope.Tx.Operations
		results := row.Result.Result.Result.Results

		for oIndex, o := range operations {
			op := NewOperation(t, &o, byte(oIndex))

			if results != nil {
				r := &(*results)[oIndex]
				AppendResult(op, r)
			}

			SerializeForBulk(op, m.buffer)
		}
	}
}

func (m *BulkMaker) makeBalancesFromMetas() {
	for tIndex, row := range m.transactionRows {
		var metas []xdr.OperationMeta

		if v1, ok := row.Meta.GetV1(); ok {
			metas = v1.Operations
		} else {
			metas, ok = row.Meta.GetOperations()
			if !ok {
				return
			}
		}

		for oIndex, e := range metas {
			pagingToken := PagingToken{
				LedgerSeq:        m.ledgerHeader.Seq,
				TransactionOrder: uint8(tIndex + 1),
				OperationOrder:   uint8(oIndex + 1),
				AuxOrder1:        1,
			}

			b := NewBalanceExtractor(e.Changes, m.ledgerHeader.CloseTime, BalanceSourceMeta, pagingToken).Extract()

			for _, balance := range b {
				SerializeForBulk(balance, m.buffer)
			}
		}
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

// 	for o := 0; o < len(metas); o++ {
// 		pagingToken := PagingToken{LedgerSeq: h.Seq, TransactionOrder: tx.Index, OperationOrder: uint8(o + 1), AuxOrder1: 1}
// 		bl := NewBalanceExtractor(metas[o].Changes, h.CloseTime, BalanceSourceMeta, pagingToken).Extract()

// 		for _, balance := range bl {
// 			SerializeForBulk(balance, b)
// 		}
// 	}
// }
