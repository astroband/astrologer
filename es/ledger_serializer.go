package es

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/astroband/astrologer/db"
	"github.com/astroband/astrologer/util"
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

// SerializeLedger serializes ledger data into ES bulk index data
func SerializeLedgerFromHistory(networkPassphrase string, meta xdr.LedgerCloseMeta, buffer *bytes.Buffer) {
	ledgerHeader := meta.V0.LedgerHeader.Header

	ledgerRow := db.LedgerHeaderRow{
		Hash:           hex.EncodeToString(meta.V0.LedgerHeader.Hash[:]),
		PrevHash:       hex.EncodeToString(ledgerHeader.PreviousLedgerHash[:]),
		BucketListHash: hex.EncodeToString(ledgerHeader.BucketListHash[:]),
		LedgerSeq:      int(ledgerHeader.LedgerSeq),
		CloseTime:      int64(ledgerHeader.ScpValue.CloseTime),
		Data:           ledgerHeader,
	}

	transactionRows := make([]db.TxHistoryRow, len(meta.V0.TxSet.Txs))
	feeRows := make([]db.TxFeeHistoryRow, len(meta.V0.TxSet.Txs))

	for i, txe := range meta.V0.TxSet.Txs {
		txHash, hashErr := util.HashTransactionInEnvelope(txe, networkPassphrase)

		if hashErr != nil {
			log.Fatalf("Failed to hash transaction #%d in ledger %d\n", i, ledgerRow.LedgerSeq)
		}

		transactionRows[i] = db.TxHistoryRow{
			ID:        hex.EncodeToString(txHash[:]),
			LedgerSeq: ledgerRow.LedgerSeq,
			Index:     i + 1,
			Envelope:  txe,
		}

		feeRows[i] = db.TxFeeHistoryRow{
			TxID:      transactionRows[i].ID,
			LedgerSeq: ledgerRow.LedgerSeq,
			Index:     i + 1,
		}

		for _, txp := range meta.V0.TxProcessing {
			if transactionRows[i].ID == hex.EncodeToString(txp.Result.TransactionHash[:]) {
				transactionRows[i].Result = txp.Result
				transactionRows[i].Meta = txp.TxApplyProcessing
				feeRows[i].Changes = txp.FeeProcessing
			}
		}
	}

	serializer := &ledgerSerializer{
		ledgerRow:       ledgerRow,
		transactionRows: transactionRows,
		feeRows:         feeRows,
		ledger:          NewLedgerHeader(&ledgerRow),
		buffer:          buffer,
	}

	serializer.serialize()
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
