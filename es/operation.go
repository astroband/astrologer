package es

import (
	"time"

	"github.com/stellar/go/xdr"
)

type AuthorizationFlag string

const (
	None                = "none"
	Full                = "full"
	MaintainLiabilities = "maintain-liabilities"
)

// Operation represents ES-serializable transaction
type Operation struct {
	TxID                 string             `json:"tx_id"`
	TxIndex              int                `json:"tx_idx"`
	Index                int                `json:"idx"`
	Seq                  int                `json:"seq"`
	PagingToken          PagingToken        `json:"paging_token"`
	CloseTime            time.Time          `json:"close_time"`
	Successful           bool               `json:"successful"`
	ResultCode           int                `json:"result_code"`
	InnerResultCode      int                `json:"inner_result_code"`
	TxSourceAccountID    string             `json:"tx_source_account_id"`
	Type                 string             `json:"type"`
	SourceAccountID      string             `json:"source_account_id,omitempty"`
	SourceAsset          *Asset             `json:"source_asset,omitempty"`
	SourceAmount         string             `json:"source_amount,omitempty"`
	AmountReceived       string             `json:"amount_received,omitempty"`
	AmountSent           string             `json:"amount_sent,omitempty"`
	DestinationAccountID string             `json:"destination_account_id,omitempty"`
	DestinationAsset     *Asset             `json:"destination_asset,omitempty"`
	DestinationAmount    string             `json:"destination_amount,omitempty"`
	OfferPrice           float64            `json:"offer_price,omitempty"`
	OfferPriceND         *Price             `json:"offer_price_n_d,omitempty"`
	OfferID              int                `json:"offer_id,omitempty"`
	TrustLimit           string             `json:"trust_limit,omitempty"`
	Authorize            AuthorizationFlag  `json:"authorize,omitempty"`
	BumpTo               int                `json:"bump_to,omitempty"`
	Path                 []*Asset           `json:"path,omitempty"`
	Thresholds           *AccountThresholds `json:"thresholds,omitempty"`
	HomeDomain           string             `json:"home_domain,omitempty"`
	InflationDest        string             `json:"inflation_dest_id,omitempty"`
	SetFlags             *AccountFlags      `json:"set_flags,omitempty"`
	ClearFlags           *AccountFlags      `json:"clear_flags,omitempty"`
	Data                 *DataEntry         `json:"data,omitempty"`
	Signer               *Signer            `json:"signer,omitempty"`

	ResultSourceAccountBalance string `json:"result_source_account_balance,omitempty"`
	ResultOffer                *Offer `json:"result_offer,omitempty"`
	ResultOfferEffect          string `json:"result_offer_effect,omitempty"`

	ResultLastAmount      string `json:"result_last_amount,omitempty"`
	ResultLastAsset       *Asset `json:"result_last_asset,omitempty"`
	ResultLastDestination string `json:"result_last_destination,omitempty"`
	ResultNoIssuer        *Asset `json:"result_no_issuer,omitempty"`

	*Memo `json:"memo,omitempty"`
}

// NewOperation creates Operation from xdr.Operation
func NewOperation(t *Transaction, o *xdr.Operation, r *[]xdr.OperationResult, n int) *Operation {
	var result *xdr.OperationResult

	if r != nil {
		result = &(*r)[n]
	}

	return ProduceOperation(t, o, result, n)
}

// DocID returns elastic document id
func (op *Operation) DocID() *string {
	s := op.PagingToken.String()
	return &s
}

// IndexName returns operations index
func (op *Operation) IndexName() IndexName {
	return opIndexName
}
