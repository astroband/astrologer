package es

import (
	"bytes"
	"log"

	"github.com/stellar/go/xdr"
)

type ledgerSerializerXDR struct {
	ledgerHeader       *LedgerHeader
	transactions       []xdr.Transaction
	transactionResults []xdr.TransactionResultPair

	buffer *bytes.Buffer
}

// SerializeLedger serializes ledger data into ES bulk index data
func SerializeLedgerFromHistory(meta xdr.LedgerCloseMeta, buffer *bytes.Buffer) {
	transactions := make([]xdr.Transaction, len(meta.V0.TxSet.Txs))
	transactionResults := make([]xdr.TransactionResultPair, len(meta.V0.TxProcessing))

	for i, txe := range meta.V0.TxSet.Txs {
		transactions[i] = txe.Tx
	}

	for i, txp := range meta.V0.TxProcessing {
		transactionResults[i] = txp.Result
	}

	serializer := &ledgerSerializerXDR{
		ledgerHeader:       NewLedgerHeaderFromHistory(meta.V0.LedgerHeader),
		transactions:       transactions,
		transactionResults: transactionResults,
		buffer:             buffer,
	}

	serializer.serialize()
}

func (s *ledgerSerializerXDR) serialize() {
	SerializeForBulk(s.ledgerHeader, s.buffer)

	for i, tx := range s.transactions {
		transaction, err := NewTransactionFromXDR(
			&transactionData{
				xdr:       tx,
				result:    s.transactionResults[i],
				index:     i + 1,
				ledgerSeq: s.ledgerHeader.Seq,
				closeTime: s.ledgerHeader.CloseTime,
			},
		)

		if err != nil {
			log.Fatal(err)
		}

		SerializeForBulk(transaction, s.buffer)

		// if transaction.Successful {
		// 	changes := s.feeRows[transaction.Index-1].Changes
		// 	s.serializeBalances(changes, transaction, nil, BalanceSourceFee)
		// }

		// 	if transaction.Successful {
		// 		changes := s.feeRows[transaction.Index-1].Changes
		// 		s.serializeBalances(changes, transaction, nil, BalanceSourceFee)
		// 	}

		// 	s.serializeOperations(transactionRow, transaction)
	}
}

// func (s *ledgerSerializerXDR) serializeOperations(transactionRow db.TxHistoryRow, transaction *Transaction) {
// 	effectsCount := 0
// 	xdrs := transactionRow.Operations()

// 	for index, xdr := range xdrs {
// 		result := transactionRow.ResultFor(index)
// 		operation := ProduceOperation(transaction, &xdr, result, index+1)
// 		SerializeForBulk(operation, s.buffer)

// 		if transaction.Successful {
// 			metas := transactionRow.MetasFor(index)
// 			if metas != nil {
// 				effectsCount = s.serializeBalances(metas.Changes, transaction, operation, BalanceSourceMeta)
// 			}

// 			s.serializeTrades(result, transaction, operation, effectsCount)

// 			h := ProduceSignerHistory(operation)
// 			if h != nil {
// 				SerializeForBulk(h, s.buffer)
// 			}
// 		}
// 	}
// }

// func (s *ledgerSerializerXDR) serializeBalances(changes xdr.LedgerEntryChanges, transaction *Transaction, operation *Operation, source BalanceSource) int {
// 	if len(changes) == 0 {
// 		return 0
// 	}

// 	pagingToken := PagingToken{
// 		LedgerSeq:        s.ledger.Seq,
// 		TransactionOrder: transaction.Index,
// 	}

// 	if operation != nil {
// 		pagingToken.OperationOrder = operation.Index
// 	}

// 	balances := ProduceBalances(changes, s.ledger.CloseTime, source, pagingToken)

// 	if len(balances) > 0 {
// 		for _, balance := range balances {
// 			SerializeForBulk(balance, s.buffer)
// 		}
// 	}

// 	return len(balances)
// }

// func (s *ledgerSerializerXDR) serializeTrades(result *xdr.OperationResult, transaction *Transaction, operation *Operation, startIndex int) int {
// 	pagingToken := PagingToken{
// 		LedgerSeq:        s.ledger.Seq,
// 		TransactionOrder: transaction.Index,
// 		OperationOrder:   operation.Index,
// 	}

// 	trades := ProduceTrades(result, operation, s.ledger.CloseTime, pagingToken, startIndex)
// 	if len(trades) > 0 {
// 		for _, trade := range trades {
// 			SerializeForBulk(&trade, s.buffer)
// 		}
// 	}

// 	return len(trades)
// }
