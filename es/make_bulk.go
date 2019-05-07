package es

import (
	"bytes"
	"fmt"

	"github.com/gzigzigzeo/stellar-core-export/db"
	"github.com/stellar/go/xdr"
)

// MakeBulk builds for bulk indexing
func MakeBulk(r db.LedgerHeaderRow, txs []db.TxHistoryRow, b *bytes.Buffer) {
	h := NewLedgerHeader(&r)
	SerializeForBulk(h, b)

	for t := 0; t < len(txs); t++ {
		var metas []xdr.OperationMeta

		txRow := &txs[t]
		ops := txRow.Envelope.Tx.Operations

		if v1, ok := txRow.Meta.GetV1(); ok {
			metas = v1.Operations
		} else {
			metas, ok = txRow.Meta.GetOperations()
		}

		tx := NewTransaction(txRow, h.CloseTime)
		SerializeForBulk(tx, b)

		for o := 0; o < len(ops); o++ {
			op := NewOperation(tx, &ops[o], byte(o))

			var r xdr.OperationResult
			if txRow.Result.Result.Result.Results != nil {
				r = (*txRow.Result.Result.Result.Results)[o]
				AppendResult(op, &r)
			}

			SerializeForBulk(op, b)
		}

		for o := 0; o < len(metas); o++ {
			id := fmt.Sprintf("%v:%v:%v", h.Seq, t, o)
			bl := ExtractBalances(metas[o].Changes, h.CloseTime, id)
			for _, balance := range bl {
				SerializeForBulk(balance, b)
			}
		}
	}
}
