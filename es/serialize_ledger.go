package es

import (
	"bytes"

	"github.com/astroband/astrologer/db"
)

// SerializeLedger serializes ledger data into ES bulk index data
func SerializeLedger(ledgerRow db.LedgerHeaderRow, transactionRows []db.TxHistoryRow, feeRows []db.TxFeeHistoryRow) {
	var b *bytes.Buffer

	ledger := NewLedgerHeader(&ledgerRow)
	SerializeForBulk(ledger, b)

	for txIndex, transactionRow := range transactionRows {
		transaction := NewTransaction(&transactionRow, ledger.CloseTime)
		SerializeForBulk(transaction, b)

		if transaction.Successful {
			fee := feeRows[txIndex]

			pagingToken := PagingToken{
				LedgerSeq:        ledger.Seq,
				TransactionOrder: txIndex + 1,
				OperationOrder:   0,
				EffectGroup:      FeeEffectPagingTokenGroup,
			}

			balances := ProduceBalances(fee.Changes, ledger.CloseTime, BalanceSourceFee, pagingToken)

			if len(balances) > 0 {
				for _, balance := range balances {
					SerializeForBulk(balance, b)
				}
			}
		}
	}
}
