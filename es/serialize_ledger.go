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
func SerializeLedger(ledgerRow db.LedgerHeaderRow, transactionRows []db.TxHistoryRow, feeRows []db.TxFeeHistoryRow) {
	ledger := NewLedgerHeader(&ledgerRow)

	serializer := &ledgerSerializer{
		ledgerRow:       ledgerRow,
		transactionRows: transactionRows,
		feeRows:         feeRows,
		ledger:          ledger,
	}

	serializer.serialize()
}

func (s *ledgerSerializer) serialize() {
	SerializeForBulk(s.ledger, s.buffer)

	for _, transactionRow := range s.transactionRows {
		transaction := NewTransaction(&transactionRow, s.ledger.CloseTime)
		SerializeForBulk(transaction, s.buffer)

		if transaction.Successful {
			changes := s.feeRows[transaction.Index].Changes
			s.serializeBalances(changes, transaction.Index, 0, BalanceSourceFee, FeeEffectPagingTokenGroup)
		}

		s.serializeOperations(transactionRow, transaction)
	}
}

func (s *ledgerSerializer) serializeOperations(transactionRow db.TxHistoryRow, transaction *Transaction) {
	xdrs := transactionRow.Operations()

	for opIndex, xdr := range xdrs {
		result := transactionRow.ResultFor(opIndex)

		operation := ProduceOperation(transaction, &xdr, result, opIndex)
		SerializeForBulk(operation, s.buffer)

		if transaction.Successful {
			// pagingToken := PagingToken{
			// 	LedgerSeq:        s.ledger.Seq,
			// 	TransactionOrder: transaction.Index + 1,
			// 	OperationOrder:   opIndex + 1,
			// }

			metas := transactionRow.MetasFor(opIndex)
			if metas != nil {
				s.serializeBalances(metas.Changes, transaction.Index, opIndex, BalanceSourceMeta, BalanceEffectPagingTokenGroup)
			}

			// extractor := NewTradeExtractor(results, op, oIndex, m.closeTime, pagingToken)
			// if extractor != nil {
			// 	trades := extractor.Extract()
			// 	for _, trade := range trades {
			// 		SerializeForBulk(&trade, m.buffer)
			// 	}
			// }

		}
	}
}

func (s *ledgerSerializer) serializeBalances(changes xdr.LedgerEntryChanges, txIndex int, opIndex int, source BalanceSource, effectGroup int) {
	if len(changes) > 0 {
		pagingToken := PagingToken{
			LedgerSeq:        s.ledger.Seq,
			TransactionOrder: txIndex + 1,
			OperationOrder:   opIndex,
			EffectGroup:      effectGroup,
		}

		balances := ProduceBalances(changes, s.ledger.CloseTime, source, pagingToken)

		if len(balances) > 0 {
			for _, balance := range balances {
				SerializeForBulk(balance, s.buffer)
			}
		}
	}
}
