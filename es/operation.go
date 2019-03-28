package es

import (
	"fmt"
	"time"

	"github.com/stellar/go/xdr"
)

type Asset struct {
	Code   string `json:"code"`
	Issuer string `json:"issuer,omitempty"`
	Key    string `json:"key"`
}

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
	SourceAsset          *Asset    `json:"source_asset,omitempty"`
	SourceAmount         int       `json:"source_amount,omitempty"`
	DestinationAccountID string    `json:"destination_account_id,omitempty"`
	DestinationAsset     *Asset    `json:"destination_asset,omitempty"`
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

	switch t := o.Body.Type; t {
	case xdr.OperationTypePayment:
		newPayment(o.Body.MustPaymentOp(), op)
	}

	return op
}

func newPayment(o xdr.PaymentOp, op *Operation) {
	op.SourceAmount = int(o.Amount)
	op.DestinationAccountID = o.Destination.Address()
	op.SourceAsset = asset(&o.Asset)
}

func asset(a *xdr.Asset) *Asset {
	var t, c, i string

	a.MustExtract(&t, &c, &i)

	if t == "native" {
		return &Asset{"native", "", "native"}
	}

	return &Asset{c, i, fmt.Sprintf("%s-%s", c, i)}
}

func (op *Operation) DocID() string {
	return op.Order
}

func (op *Operation) IndexName() string {
	return opIndexName
}
