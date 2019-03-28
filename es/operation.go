package es

import (
	"fmt"
	"time"

	"github.com/stellar/go/xdr"
)

// Transaction represents ES-serializable transaction
type Operation struct {
	TxID                 string    `json:"tx_id"`
	TxIndex              byte      `json:"tx_idx"`
	Index                byte      `json:"idx"`
	Seq                  int       `json:"seq"`
	Order                string    `json:"order"`
	CloseTime            time.Time `json:"close_time"`
	Succesful            bool      `json:"successful"`
	ResultCode           int       `json:"result_code"`
	TxSourceAccountID    string    `json:"tx_source_account_id"`
	Type                 string    `json:"type"`
	SourceAccountID      string    `json:"source_account_id,omitempty"`
	SourceAsset          string    `json:"source_asset,omitempty"`
	SourceAmount         int       `json:"source_amount,omitempty"`
	DestinationAccountID string    `json:"destination_account_id,omitempty"`
	DestinationAsset     string    `json:"destination_asset,omitempty"`
	DestinationAmount    int       `json:"destination_amount,omitempty"`
	StartingBalance      int       `json:"starting_balance,omitempty"`
	OfferPrice           int       `json:"offer_price,omitempty"`
	OfferID              int       `json:"offer_id,omitempty"`
	TrustLimit           int       `json:"trust_limit,omitempty"`
	Authorize            bool      `json:"authorize,omitempty"`
	BumpTo               int       `json:"bump_to,omitempty"`

	*Memo `json:"memo,omitempty"`
}

// NewOperation creates LedgerHeader from LedgerHeaderRow
func NewOperation(t *Transaction, o *xdr.Operation, n byte) *Operation {
	sourceAccountID := t.SourceAccountID
	if o.SourceAccount != nil {
		sourceAccountID = o.SourceAccount.Address()
	}

	op := &Operation{
		TxID:              t.ID,
		TxIndex:           t.Index,
		Index:             n,
		Seq:               t.Seq,
		Order:             fmt.Sprintf("%s:%d", t.Order, n),
		CloseTime:         t.CloseTime,
		Succesful:         true, // TODO: Implement
		ResultCode:        0,    // TODO: Implement
		TxSourceAccountID: t.SourceAccountID,
		Type:              o.Body.Type.String(),
		SourceAccountID:   sourceAccountID,

		Memo: t.Memo,
	}

	return op
}

func (op *Operation) DocID() string {
	return op.Order
}

func (op *Operation) IndexName() string {
	return opIndexName
}
