package es

import (
	"bytes"
	"time"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
)

var (
	lastOperationIndex       uint8 = 255
	balanceFromMetaAux1Order uint8 = 1
	balanceFromFeeAux1Order  uint8 = 2
)

// BulkMaker creates es bulk from ledger data
type BulkMaker struct {
	ledgerRow       db.LedgerHeaderRow
	ledgerHeader    *LedgerHeader
	seq             int
	closeTime       time.Time
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
		seq:             h.Seq,
		closeTime:       h.CloseTime,
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
	m.makeOperationsWithResultsAndTrades()
	m.makeBalancesFromMetas()
	m.makeBalancesFromFeeHistory()
}

func (m *BulkMaker) makeLedger() {
	SerializeForBulk(m.ledgerHeader, m.buffer)
}

func (m *BulkMaker) makeTransactions() {
	for _, transaction := range m.transactions {
		SerializeForBulk(transaction, m.buffer)
	}
}

func (m *BulkMaker) makeOperationsWithResultsAndTrades() {
	for tIndex, t := range m.transactions {
		row := m.transactionRows[tIndex]
		operations := row.Envelope.Tx.Operations
		results := row.Result.Result.Result.Results

		for oIndex, o := range operations {
			op := NewOperation(t, &o, results, byte(oIndex))
			SerializeForBulk(op, m.buffer)

			pagingToken := PagingToken{
				LedgerSeq:        m.seq,
				TransactionOrder: uint8(tIndex + 1),
				OperationOrder:   uint8(oIndex + 1),
			}

			extractor := NewTradeExtractor(results, op, oIndex, m.closeTime, pagingToken)
			if extractor != nil {
				trades := extractor.Extract()
				for _, trade := range trades {
					SerializeForBulk(&trade, m.buffer)
				}
			}
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
				LedgerSeq:        m.seq,
				TransactionOrder: uint8(tIndex + 1),
				OperationOrder:   uint8(oIndex + 1),
				AuxOrder1:        balanceFromMetaAux1Order,
			}
			b := NewBalanceExtractor(e.Changes, m.closeTime, BalanceSourceMeta, pagingToken).Extract()

			for _, balance := range b {
				SerializeForBulk(balance, m.buffer)
			}
		}
	}
}

func (m *BulkMaker) makeBalancesFromFeeHistory() {
	for tIndex, fee := range m.fees {
		pagingToken := PagingToken{
			LedgerSeq:        m.seq,
			TransactionOrder: uint8(tIndex + 1),
			OperationOrder:   lastOperationIndex,
			AuxOrder1:        balanceFromFeeAux1Order,
		}

		bl := NewBalanceExtractor(fee.Changes, m.closeTime, BalanceSourceFee, pagingToken).Extract()

		for _, balance := range bl {
			SerializeForBulk(balance, m.buffer)
		}
	}
}
