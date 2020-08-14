package db

import (
	"fmt"
	"github.com/stellar/go/xdr"
	"log"
)

// TxHistoryRow represents row of txhistory table
type TxHistoryRow struct {
	ID        string                    `db:"txid"`
	LedgerSeq int                       `db:"ledgerseq"`
	Index     int                       `db:"txindex"`
	Envelope  xdr.TransactionEnvelope   `db:"txbody"`
	Result    xdr.TransactionResultPair `db:"txresult"`
	Meta      xdr.TransactionMeta       `db:"txmeta"`
}

// TxHistoryRowForSeq returns transactions for specified ledger sorted by index
func (db *Client) TxHistoryRowForSeq(seq int) []TxHistoryRow {
	txs := []TxHistoryRow{}

	err := db.rawClient.Select(&txs, "SELECT * FROM txhistory WHERE ledgerseq = $1 ORDER BY txindex", seq)
	if err != nil {
		log.Fatal(err)
	}

	return txs
}

// Operations returns operations array
func (tx *TxHistoryRow) Operations() (error, []xdr.Operation) {
	switch tx.Envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		return nil, tx.Envelope.V0.Tx.Operations
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		return nil, tx.Envelope.V1.Tx.Operations
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		return nil, tx.Envelope.FeeBump.Tx.InnerTx.V1.Tx.Operations
	default:
		return fmt.Errorf("Unknown tx envelope type %s", tx.Envelope.Type), make([]xdr.Operation, 0)
	}
}

// ResultFor returns result for operation index
func (tx *TxHistoryRow) ResultFor(index int) (result *xdr.OperationResult) {
	results := tx.Result.Result.Result.Results

	if results != nil {
		result = &(*results)[index]
	}

	return result
}

// MetasFor returns meta for operation index
func (tx *TxHistoryRow) MetasFor(index int) (result *xdr.OperationMeta) {
	if v1, ok := tx.Meta.GetV1(); ok {
		ops := v1.Operations
		return &ops[index]
	}

	ops, ok := tx.Meta.GetOperations()
	if !ok {
		return nil
	}

	return &ops[index]
}
