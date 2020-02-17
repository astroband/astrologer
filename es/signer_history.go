package es

import (
	"time"
)

// SignerHistory represents signer change entry, still the question if it is required
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

// ProduceSignerHistory creates new signer history entry
func ProduceSignerHistory(o *Operation) (h *SignerHistory) {
	s := o.Signer

	if s == nil {
		return nil
	}

	token := PagingToken{
		LedgerSeq:        o.Seq,
		TransactionOrder: o.TxIndex,
		OperationOrder:   o.Index,
	}

	entry := &SignerHistory{
		PagingToken:     token,
		AccountID:       o.SourceAccountID,
		Signer:          o.Signer.ID,
		Type:            o.Signer.Type,
		Weight:          o.Signer.Weight,
		TxIndex:         o.TxIndex,
		Index:           o.Index,
		Seq:             o.Seq,
		LedgerCloseTime: o.CloseTime,
	}

	return entry
}

// DocID balance es document id
func (t *SignerHistory) DocID() *string {
	s := t.PagingToken.String()
	return &s
}

// IndexName balances index name
func (t *SignerHistory) IndexName() IndexName {
	return signerHistoryIndexName
}
