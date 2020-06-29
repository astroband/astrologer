package es

import (
	"bytes"
	"fmt"

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
func SerializeLedger(ledgerRow db.LedgerHeaderRow, transactionRows []db.TxHistoryRow, feeRows []db.TxFeeHistoryRow, buffer *bytes.Buffer) error {
	ledger := NewLedgerHeader(&ledgerRow)

	serializer := &ledgerSerializer{
		ledgerRow:       ledgerRow,
		transactionRows: transactionRows,
		feeRows:         feeRows,
		ledger:          ledger,
		buffer:          buffer,
	}

	return serializer.serialize()
}

func (s *ledgerSerializer) serialize() error {
	SerializeForBulk(s.ledger, s.buffer)

	for _, transactionRow := range s.transactionRows {
		transaction, err := s.NewTransaction(&transactionRow, s.ledger.CloseTime)

		if err != nil {
			return err
		}

		SerializeForBulk(transaction, s.buffer)

		if transaction.Successful {
			changes := s.feeRows[transaction.Index-1].Changes
			s.serializeBalances(changes, transaction, nil, BalanceSourceFee)
		}

		s.serializeOperations(transactionRow, transaction)
	}

	return nil
}

func (s *ledgerSerializer) serializeOperations(transactionRow db.TxHistoryRow, transaction *Transaction) error {
	effectsCount := 0
	err, xdrs := transactionRow.Operations()

	if err != nil {
		return err
	}

	for index, xdr := range xdrs {
		result := transactionRow.ResultFor(index)
		operation, err := ProduceOperation(transaction, &xdr, result, index+1)

		if err != nil {
			return fmt.Errorf("Failed to serialize operation with index %d in tx %s: %w", index, transaction.ID, err)
		}

		SerializeForBulk(operation, s.buffer)

		if transaction.Successful {
			metas := transactionRow.MetasFor(index)
			if metas != nil {
				effectsCount = s.serializeBalances(metas.Changes, transaction, operation, BalanceSourceMeta)
			}

			s.serializeTrades(result, transaction, operation, effectsCount)

			h := ProduceSignerHistory(operation)
			if h != nil {
				SerializeForBulk(h, s.buffer)
			}
		}
	}

	return nil
}

func (s *ledgerSerializer) serializeBalances(changes xdr.LedgerEntryChanges, transaction *Transaction, operation *Operation, source BalanceSource) int {
	if len(changes) == 0 {
		return 0
	}

	pagingToken := PagingToken{
		LedgerSeq:        s.ledger.Seq,
		TransactionOrder: transaction.Index,
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

	return len(balances)
}

func (s *ledgerSerializer) serializeTrades(result *xdr.OperationResult, transaction *Transaction, operation *Operation, startIndex int) int {
	pagingToken := PagingToken{
		LedgerSeq:        s.ledger.Seq,
		TransactionOrder: transaction.Index,
		OperationOrder:   operation.Index,
	}

	trades := ProduceTrades(result, operation, s.ledger.CloseTime, pagingToken, startIndex)
	if len(trades) > 0 {
		for _, trade := range trades {
			SerializeForBulk(&trade, s.buffer)
		}
	}

	return len(trades)
}
