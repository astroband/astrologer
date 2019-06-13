package es

import (
	"time"

	"github.com/stellar/go/xdr"
)

// Operation represents ES-serializable transaction
type Operation struct {
	TxID                 string             `json:"tx_id"`
	TxIndex              int                `json:"tx_idx"`
	Index                int                `json:"idx"`
	Seq                  int                `json:"seq"`
	PagingToken          PagingToken        `json:"paging_token"`
	CloseTime            time.Time          `json:"close_time"`
	Succesful            bool               `json:"successful"`
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
	Authorize            bool               `json:"authorize,omitempty"`
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
	op := ProduceOperation(t, o, r, n)

	if r != nil {
		result := &(*r)[n]
		AppendResult(op, result)
	}

	return op

	// 	sourceAccountID := t.SourceAccountID

	// 	if o.SourceAccount != nil {
	// 		sourceAccountID = o.SourceAccount.Address()
	// 	}

	// 	op := &Operation{
	// 		TxID:              t.ID,
	// 		TxIndex:           t.Index,
	// 		Index:             n,
	// 		Seq:               t.Seq,
	// 		PagingToken:       PagingToken{LedgerSeq: t.Seq, TransactionOrder: t.Index + 1, OperationOrder: n + 1},
	// 		CloseTime:         t.CloseTime,
	// 		TxSourceAccountID: t.SourceAccountID,
	// 		Type:              o.Body.Type.String(),
	// 		SourceAccountID:   sourceAccountID,

	// 		Memo: t.Memo,
	// 	}

	// 	switch t := o.Body.Type; t {
	// 	case xdr.OperationTypeCreateAccount:
	// 		newCreateAccount(o.Body.MustCreateAccountOp(), op)
	// 	case xdr.OperationTypePayment:
	// 		newPayment(o.Body.MustPaymentOp(), op)
	// 	case xdr.OperationTypePathPayment:
	// 		newPathPayment(o.Body.MustPathPaymentOp(), op)
	// 	case xdr.OperationTypeManageSellOffer:
	// 		newManageSellOffer(o.Body.MustManageSellOfferOp(), op)
	// 	case xdr.OperationTypeManageBuyOffer:
	// 		newManageBuyOffer(o.Body.MustManageBuyOfferOp(), op)
	// 	case xdr.OperationTypeCreatePassiveSellOffer:
	// 		newCreatePassiveSellOffer(o.Body.MustCreatePassiveSellOfferOp(), op)
	// 	case xdr.OperationTypeSetOptions:
	// 		newSetOptions(o.Body.MustSetOptionsOp(), op)
	// 	case xdr.OperationTypeChangeTrust:
	// 		newChangeTrust(o.Body.MustChangeTrustOp(), op)
	// 	case xdr.OperationTypeAllowTrust:
	// 		newAllowTrust(o.Body.MustAllowTrustOp(), op)
	// 	case xdr.OperationTypeAccountMerge:
	// 		newAccountMerge(o.Body.MustDestination(), op)
	// 	case xdr.OperationTypeManageData:
	// 		newManageData(o.Body.MustManageDataOp(), op)
	// 	case xdr.OperationTypeBumpSequence:
	// 		newBumpSequence(o.Body.MustBumpSequenceOp(), op)
	// 	}

	// 	if r != nil {
	// 		result := &(*r)[n]
	// 		AppendResult(op, result)
	// 	}

	// 	return op
	// }

	// func newCreateAccount(o xdr.CreateAccountOp, op *Operation) {
	// 	op.SourceAmount = amount.String(o.StartingBalance)
	// 	op.DestinationAccountID = o.Destination.Address()
	// }

	// func newPayment(o xdr.PaymentOp, op *Operation) {
	// 	op.SourceAmount = amount.String(o.Amount)
	// 	op.DestinationAccountID = o.Destination.Address()
	// 	op.SourceAsset = NewAsset(&o.Asset)
	// }

	// func newPathPayment(o xdr.PathPaymentOp, op *Operation) {
	// 	op.DestinationAccountID = o.Destination.Address()
	// 	op.DestinationAmount = amount.String(o.DestAmount)
	// 	op.DestinationAsset = NewAsset(&o.DestAsset)

	// 	op.SourceAmount = amount.String(o.SendMax)
	// 	op.SourceAsset = NewAsset(&o.SendAsset)

	// 	op.Path = make([]*Asset, len(o.Path))

	// 	for i, a := range o.Path {
	// 		op.Path[i] = NewAsset(&a)
	// 	}
	// }

	// func newManageSellOffer(o xdr.ManageSellOfferOp, op *Operation) {
	// 	op.SourceAmount = amount.String(o.Amount)
	// 	op.SourceAsset = NewAsset(&o.Buying)
	// 	op.OfferID = int(o.OfferId)
	// 	op.OfferPrice, _ = big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()
	// 	op.OfferPriceND = &Price{int(o.Price.N), int(o.Price.D)}
	// 	op.DestinationAsset = NewAsset(&o.Selling)
	// }

	// func newManageBuyOffer(o xdr.ManageBuyOfferOp, op *Operation) {
	// 	op.SourceAmount = amount.String(o.BuyAmount)
	// 	op.SourceAsset = NewAsset(&o.Selling)
	// 	op.DestinationAsset = NewAsset(&o.Buying)
	// 	op.OfferID = int(o.OfferId)
	// 	op.OfferPrice, _ = big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()
	// 	op.OfferPriceND = &Price{int(o.Price.N), int(o.Price.D)}
	// }

	// func newCreatePassiveSellOffer(o xdr.CreatePassiveSellOfferOp, op *Operation) {
	// 	op.SourceAmount = amount.String(o.Amount)
	// 	op.SourceAsset = NewAsset(&o.Buying)
	// 	op.OfferPrice, _ = big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()
	// 	op.OfferPriceND = &Price{int(o.Price.N), int(o.Price.D)}
	// 	op.DestinationAsset = NewAsset(&o.Selling)
	// }

	// func newSetOptions(o xdr.SetOptionsOp, op *Operation) {
	// 	if o.InflationDest != nil {
	// 		op.InflationDest = o.InflationDest.Address()
	// 	}

	// 	if o.HomeDomain != nil {
	// 		op.HomeDomain = string(*o.HomeDomain)
	// 	}

	// 	if (o.LowThreshold != nil) || (o.MedThreshold != nil) || (o.HighThreshold != nil) || (o.MasterWeight != nil) {
	// 		op.Thresholds = &AccountThresholds{}

	// 		if o.LowThreshold != nil {
	// 			op.Thresholds.Low = new(byte)
	// 			*op.Thresholds.Low = byte(*o.LowThreshold)
	// 		}

	// 		if o.MedThreshold != nil {
	// 			op.Thresholds.Medium = new(byte)
	// 			*op.Thresholds.Medium = byte(*o.MedThreshold)
	// 		}

	// 		if o.HighThreshold != nil {
	// 			op.Thresholds.High = new(byte)
	// 			*op.Thresholds.High = byte(*o.HighThreshold)
	// 		}

	// 		if o.MasterWeight != nil {
	// 			op.Thresholds.Master = new(byte)
	// 			*op.Thresholds.Master = byte(*o.MasterWeight)
	// 		}
	// 	}

	// 	if o.SetFlags != nil {
	// 		op.SetFlags = flags(int(*o.SetFlags))
	// 	}

	// 	if o.ClearFlags != nil {
	// 		op.ClearFlags = flags(int(*o.ClearFlags))
	// 	}

	// 	if o.Signer != nil {
	// 		op.Signer = &Signer{
	// 			o.Signer.Key.Address(),
	// 			int(o.Signer.Weight),
	// 		}
	// 	}
	// }

	// func newChangeTrust(o xdr.ChangeTrustOp, op *Operation) {
	// 	op.DestinationAmount = amount.String(o.Limit)
	// 	op.DestinationAsset = NewAsset(&o.Line)
	// }

	// func newAllowTrust(o xdr.AllowTrustOp, op *Operation) {
	// 	a := o.Asset.ToAsset(o.Trustor)

	// 	op.DestinationAsset = NewAsset(&a)
	// 	op.DestinationAccountID = o.Trustor.Address()
	// 	op.Authorize = o.Authorize
	// }

	// func newAccountMerge(d xdr.AccountId, op *Operation) {
	// 	op.DestinationAccountID = d.Address()
	// }

	// func newBumpSequence(o xdr.BumpSequenceOp, op *Operation) {
	// 	op.BumpTo = int(o.BumpTo)
	// }

	// // TODO: Apply some magic to the value
	// func newManageData(o xdr.ManageDataOp, op *Operation) {
	// 	op.Data = &DataEntry{Name: string(o.DataName)}
	// 	if o.DataValue != nil {
	// 		op.Data.Value = string(*o.DataValue)
	// 	}
	// }

	// func flags(f int) *AccountFlags {
	// 	l := xdr.AccountFlags(f)

	// 	return &AccountFlags{
	// 		l&xdr.AccountFlagsAuthRequiredFlag != 0,
	// 		l&xdr.AccountFlagsAuthRevocableFlag != 0,
	// 		l&xdr.AccountFlagsAuthImmutableFlag != 0,
	// 	}
	// }
}

// DocID returns elastic document id
func (op *Operation) DocID() *string {
	s := op.PagingToken.String()
	return &s
}

// IndexName returns operations index
func (op *Operation) IndexName() string {
	return opIndexName
}
