package es

import (
	"github.com/stellar/go/xdr"
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

// ProduceSignerHistoryFromTxMeta handles pre-authorized transactions
//
// Stellar allows user to create pre-authorized transactions, which
// can be submitted later by other users. This mechanism is implemented
// using multi-sig. When such pre-authorized transaction is executed, the "signer"
// for this transaction is deleted. We should ingest this change too
//
// You can read about pre-authorized transactions process here:
// https://medium.com/@katopz/stellar-pre-signed-transaction-d93e91191c15
func ProduceSignerHistoryFromTxMeta(tx *Transaction) (h *SignerHistory) {
	currentSigners := make(map[string][]xdr.Signer)

	for _, change := range tx.meta.TxChanges {
		if change.EntryType() != xdr.LedgerEntryTypeAccount {
			continue
		}

		switch t := change.Type; t {
		case xdr.LedgerEntryChangeTypeLedgerEntryState:
			accountData := change.MustState().Data.MustAccount()
			accountId := accountData.AccountId.Address()

			currentSigners[accountId] = accountData.Signers

		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			accountData := change.MustUpdated().Data.MustAccount()
			accountId := accountData.AccountId.Address()

			if len(currentSigners[accountId]) == len(accountData.Signers) {
				continue
			}

			for _, signer := range currentSigners[accountId] {
				if contains(accountData.Signers, signer) {
					continue
				}

				s := NewSigner(&signer)

				token := PagingToken{
					LedgerSeq:        tx.Seq,
					TransactionOrder: tx.Index,
				}

				entry := &SignerHistory{
					PagingToken:     token,
					AccountID:       accountId,
					Signer:          s.ID,
					Type:            s.Type,
					Weight:          0,
					TxIndex:         tx.Index,
					Index:           0,
					Seq:             tx.Seq,
					LedgerCloseTime: tx.CloseTime,
				}

				// we can return quickly, because there can be only one
				// signer change in pre-auth tx
				return entry
			}
		}
	}

	return nil
}

// contains checks, whether given array of signers contains given
// particular signer or nor
func contains(array []xdr.Signer, signer xdr.Signer) bool {
	for _, s := range array {
		if s.Key.Address() == signer.Key.Address() && s.Weight == signer.Weight {
			return true
		}
	}

	return false
}

// ProduceSignerHistoryFromOperation creates new signer history entry
func ProduceSignerHistoryFromOperation(o *Operation) (h *SignerHistory) {
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
