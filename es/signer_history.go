package es

import (
	"time"
)

// SignerHistory represents trade entry
type SignerHistory struct {
	PagingToken     PagingToken `json:"paging_token"`
	AccountID       string      `json:"account_id"`
	Signer          string      `json:"signer"`
	Type            int         `json:"type"`
	Weight          int         `json:"weight"`
	TxIndex         int         `json:"tx_idx"`
	Index           int         `json:"idx"`
	Seq             int         `json:"seq"`
	LedgerCloseTime time.Time   `json:"ledger_close_time"`
}

// DocID balance es document id
func (t *SignerHistory) DocID() *string {
	s := t.PagingToken.String()
	return &s
}

// IndexName balances index name
func (t *SignerHistory) IndexName() string {
	return signerHistoryIndexName
}
