package es

import (
	"bytes"

	"github.com/astroband/astrologer/db"
	"github.com/stellar/go/xdr"
)

// MakeBulk builds for bulk indexing
func MakeBulk(r db.LedgerHeaderRow, txs []db.TxHistoryRow, fees []db.TxFeeHistoryRow, b *bytes.Buffer) {
	h := NewLedgerHeader(&r)
	SerializeForBulk(h, b)

	for t := 0; t < len(txs); t++ {
		var metas []xdr.OperationMeta

		txRow := &txs[t]
		ops := txRow.Envelope.Tx.Operations
		results := txRow.Result.Result.Result.Results

		if v1, ok := txRow.Meta.GetV1(); ok {
			metas = v1.Operations
		} else {
			metas, ok = txRow.Meta.GetOperations()
		}

		tx := NewTransaction(txRow, h.CloseTime)
		SerializeForBulk(tx, b)

		for o := 0; o < len(ops); o++ {
			op := NewOperation(tx, &ops[o], byte(o))

			if results != nil {
				AppendResult(op, &(*results)[o])
			}

			SerializeForBulk(op, b)
		}

		for o := 0; o < len(metas); o++ {
			order := Order{LedgerSeq: h.Seq, TransactionOrder: tx.Index, OperationOrder: uint8(o), AuxOrder1: 0}
			bl := NewBalanceExtractor(metas[o].Changes, h.CloseTime, BalanceSourceMeta, order).Extract()

			for _, balance := range bl {
				SerializeForBulk(balance, b)
			}
		}
	}

	// for o := 0; o < len(fees); o++ {
	// 	fee := fees[o]

	// 	order := Order{LedgerSeq: h.Seq, TransactionOrder: tx.Index, OperationOrder: o, AuxOrder1: 1}
	// 	bl := NewBalanceExtractor(fee.Changes, h.CloseTime, BalanceSourceFee, order).Extract()

	// 	for _, balance := range bl {
	// 		SerializeForBulk(balance, b)
	// 	}
	// }
}
