package es

import (
	"bytes"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
)

type ledgerSerializer struct {
	ledgerRow       db.LedgerHeaderRow
	transactionRows []db.TxHistoryRow
	feeRows         []db.TxFeeHistoryRow
	ledger          *LedgerHeader

	buffer *bytes.Buffer
}

// SerializeLedger serializes ledger data into ES bulk index data
func SerializeLedger(ledgerRow db.LedgerHeaderRow, transactionRows []db.TxHistoryRow, feeRows []db.TxFeeHistoryRow, buffer *bytes.Buffer) {
	ledger := NewLedgerHeader(&ledgerRow)

	serializer := &ledgerSerializer{
		ledgerRow:       ledgerRow,
		transactionRows: transactionRows,
		feeRows:         feeRows,
		ledger:          ledger,
		buffer:          buffer,
	}

	serializer.serialize()
}

func (s *ledgerSerializer) serialize() {
	SerializeForBulk(s.ledger, s.buffer)

	for i, transactionRow := range s.transactionRows {
		transaction := NewTransaction(&transactionRow, &s.feeRows[i], s.ledger.CloseTime)
		SerializeForBulk(transaction, s.buffer)

		if transaction.Successful {
			changes := s.feeRows[transaction.Index-1].Changes
			s.serializeBalances(changes, transaction, nil, BalanceSourceFee, FeeEffectPagingTokenGroup)
			s.serializeEvents(
				s.transactionRows[transaction.Index-1].Meta.V1.Operations,
				transaction.Index,
				nil,
				OfferEffectPagingTokenGroup,
			)
		}

		s.serializeOperations(transactionRow, transaction)
	}
}

func (s *ledgerSerializer) serializeOperations(transactionRow db.TxHistoryRow, transaction *Transaction) {
	xdrs := transactionRow.Operations()

	for index, xdr := range xdrs {
		result := transactionRow.ResultFor(index)
		operation := ProduceOperation(transaction, &xdr, result, index+1)
		SerializeForBulk(operation, s.buffer)

		if transaction.Successful {
			metas := transactionRow.MetasFor(index)
			if metas != nil {
				s.serializeBalances(metas.Changes, transaction, operation, BalanceSourceMeta, BalanceEffectPagingTokenGroup)
			}

			s.serializeTrades(result, transaction, operation)

			h := ProduceSignerHistory(operation)
			if h != nil {
				SerializeForBulk(h, s.buffer)
			}
		}
	}
}

func (s *ledgerSerializer) serializeBalances(changes xdr.LedgerEntryChanges, transaction *Transaction, operation *Operation, source BalanceSource, effectGroup int) {
	if len(changes) == 0 {
		return
	}

	pagingToken := PagingToken{
		LedgerSeq:        s.ledger.Seq,
		TransactionOrder: transaction.Index,
		EffectGroup:      effectGroup,
	}

	if operation != nil {
		pagingToken.OperationOrder = operation.Index
	}

	balances := ProduceBalances(changes, s.ledger.CloseTime, source, pagingToken)

	if len(balances) > 0 {
		for _, balance := range balances {
			SerializeForBulk(balance, s.buffer)
		}
	}
}

func (s *ledgerSerializer) serializeEvents(changes []xdr.OperationMeta, txIndex int, opIndex *int, effectGroup int) {
	if len(changes) == 0 {
		return
	}

	pagingToken := PagingToken{
		LedgerSeq:        s.ledger.Seq,
		TransactionOrder: txIndex,
		EffectGroup:      effectGroup,
	}

	if opIndex != nil {
		pagingToken.OperationOrder = *opIndex
	}

	events := ProduceEvents(changes, s.ledger.Seq, s.ledger.CloseTime, pagingToken)

	if len(events) > 0 {
		for _, event := range events {
			SerializeForBulk(event, s.buffer)
		}
	}
}

func (s *ledgerSerializer) serializeTrades(result *xdr.OperationResult, transaction *Transaction, operation *Operation) {
	pagingToken := PagingToken{
		LedgerSeq:        s.ledger.Seq,
		TransactionOrder: transaction.Index,
		OperationOrder:   operation.Index,
	}

	trades := ProduceTrades(result, operation, s.ledger.CloseTime, pagingToken)
	if len(trades) > 0 {
		for _, trade := range trades {
			SerializeForBulk(&trade, s.buffer)
		}
	}
}
