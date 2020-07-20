package es

import (
	"bytes"
	"log"

	"github.com/astroband/astrologer/support"
	"github.com/stellar/go/xdr"
)

type txVersion int

const (
	v0 = iota
	v1 = iota
)

type txSetMember struct {
	v     txVersion
	xdrV0 *xdr.TransactionV0
	xdrV1 *xdr.Transaction
}

type ledgerSerializerXDR struct {
	ledgerHeader       *LedgerHeader
	txSet              []txSetMember
	transactionResults []xdr.TransactionResultPair
	transactionMetas   []xdr.TransactionMeta
	changes            []xdr.LedgerEntryChanges

	buffer *bytes.Buffer
}

// SerializeLedger serializes ledger data into ES bulk index data
func SerializeLedgerFromHistory(meta xdr.LedgerCloseMeta, buffer *bytes.Buffer) {
	serializer := &ledgerSerializerXDR{
		txSet:              make([]txSetMember, len(meta.V0.TxSet.Txs)),
		ledgerHeader:       NewLedgerHeaderFromHistory(meta.V0.LedgerHeader),
		transactionResults: make([]xdr.TransactionResultPair, len(meta.V0.TxProcessing)),
		changes:            make([]xdr.LedgerEntryChanges, len(meta.V0.TxProcessing)),
		transactionMetas:   make([]xdr.TransactionMeta, len(meta.V0.TxProcessing)),
		buffer:             buffer,
	}

	for i, txe := range meta.V0.TxSet.Txs {
		switch txe.Type {
		case xdr.EnvelopeTypeEnvelopeTypeTxV0:
			serializer.txSet[i] = txSetMember{v: v0, xdrV0: &txe.V0.Tx}
		case xdr.EnvelopeTypeEnvelopeTypeTx:
			serializer.txSet[i] = txSetMember{v: v1, xdrV1: &txe.V1.Tx}
		default:
			log.Fatal("Unknown type")
		}
	}

	for i, txp := range meta.V0.TxProcessing {
		serializer.transactionResults[i] = txp.Result
		serializer.changes[i] = txp.FeeProcessing
		serializer.transactionMetas[i] = txp.TxApplyProcessing
	}

	serializer.serialize()
}

func (s *ledgerSerializerXDR) serialize() {
	SerializeForBulk(s.ledgerHeader, s.buffer)

	for i, tx := range s.txSet {
		txData := transactionData{
			result:    s.transactionResults[i],
			index:     i + 1,
			ledgerSeq: s.ledgerHeader.Seq,
			closeTime: s.ledgerHeader.CloseTime,
		}

		switch tx.v {
		case v0:
			txData.v = v0
			txData.xdrV0 = tx.xdrV0
		case v1:
			txData.v = v1
			txData.xdrV1 = tx.xdrV1
		}

		transaction, err := NewTransactionFromXDR(&txData)

		if err != nil {
			log.Fatal(err)
		}

		SerializeForBulk(transaction, s.buffer)

		if transaction.Successful {
			changes := s.changes[i]
			s.serializeBalances(changes, transaction, nil, BalanceSourceFee)
		}

		// 	s.serializeOperations(transactionRow, transaction)
	}
}

func (s *ledgerSerializerXDR) serializeOperations(operations []xdr.Operation, operationResults *[]xdr.OperationResult, transaction *Transaction) error {
	// effectsCount := 0

	for index, operation := range operations {
		result := (*operationResults)[index]
		operation, err := ProduceOperation(transaction, &operation, &result, index+1)

		if err != nil {
			return err
		}

		SerializeForBulk(operation, s.buffer)

		if transaction.Successful {
			metas := support.OperationMeta(s.transactionMetas[transaction.Index], index)
			if metas != nil {
				// effectsCount = s.serializeBalances(metas.Changes, transaction, operation, BalanceSourceMeta)
				s.serializeBalances(metas.Changes, transaction, operation, BalanceSourceMeta)
			}

			// s.serializeTrades(result, transaction, operation, effectsCount)

			h := ProduceSignerHistory(operation)
			if h != nil {
				SerializeForBulk(h, s.buffer)
			}
		}
	}

	return nil
}

func (s *ledgerSerializerXDR) serializeBalances(changes xdr.LedgerEntryChanges, transaction *Transaction, operation *Operation, source BalanceSource) int {
	if len(changes) == 0 {
		return 0
	}

	pagingToken := PagingToken{
		LedgerSeq:        s.ledgerHeader.Seq,
		TransactionOrder: transaction.Index,
	}

	if operation != nil {
		pagingToken.OperationOrder = operation.Index
	}

	balances := ProduceBalances(changes, s.ledgerHeader.CloseTime, source, pagingToken)

	if len(balances) > 0 {
		for _, balance := range balances {
			SerializeForBulk(balance, s.buffer)
		}
	}

	return len(balances)
}

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
