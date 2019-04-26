package db

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/guregu/null"
	"github.com/gzigzigzeo/stellar-core-export/config"
	"github.com/stellar/go/xdr"
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
func TxHistoryRowForSeq(seq int) []TxHistoryRow {
	txs := []TxHistoryRow{}

	err := config.DB.Select(&txs, "SELECT * FROM txhistory WHERE ledgerseq = $1 ORDER BY txindex", seq)
	if err != nil {
		log.Fatal(err)
	}

	return txs
}

// MemoValue Returns clean memo value, this is copy paste from horizon internal package
func (tx *TxHistoryRow) MemoValue() null.String {
	var (
		value string
		valid bool
	)
	switch tx.Envelope.Tx.Memo.Type {
	case xdr.MemoTypeMemoNone:
		value, valid = "", false
	case xdr.MemoTypeMemoText:
		scrubbed := utf8Scrub(tx.Envelope.Tx.Memo.MustText())
		notnull := strings.Join(strings.Split(scrubbed, "\x00"), "")
		value, valid = notnull, true
	case xdr.MemoTypeMemoId:
		value, valid = fmt.Sprintf("%d", tx.Envelope.Tx.Memo.MustId()), true
	case xdr.MemoTypeMemoHash:
		hash := tx.Envelope.Tx.Memo.MustHash()
		value, valid =
			base64.StdEncoding.EncodeToString(hash[:]),
			true
	case xdr.MemoTypeMemoReturn:
		hash := tx.Envelope.Tx.Memo.MustRetHash()
		value, valid =
			base64.StdEncoding.EncodeToString(hash[:]),
			true
	default:
		panic(fmt.Errorf("invalid memo type: %v", tx.Envelope.Tx.Memo.Type))
	}

	return null.NewString(value, valid)
}
