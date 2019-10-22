package es

import (
	"math/big"
	"strings"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// AppendResult Appends operation result
func (f *operationFactory) assignResult() {
	if f.result == nil {
		return
	}

	r := f.result

	f.operation.ResultCode = int(r.Code)

	if r.Code == xdr.OperationResultCodeOpInner {
		switch t := r.Tr.Type; t {
		case xdr.OperationTypeCreateAccount:
			f.assignCreateAccountResult(r.Tr.MustCreateAccountResult())
		case xdr.OperationTypePayment:
			f.assignPaymentResult(r.Tr.MustPaymentResult())
		case xdr.OperationTypePathPaymentStrictReceive:
			f.assignPathPaymentStrictReceiveResult(r.Tr.MustPathPaymentStrictReceiveResult())
		case xdr.OperationTypePathPaymentStrictSend:
			f.assignPathPaymentStrictSendResult(r.Tr.MustPathPaymentStrictSendResult())
		case xdr.OperationTypeManageSellOffer:
			f.assignManageSellOfferResult(r.Tr.MustManageSellOfferResult())
		case xdr.OperationTypeManageBuyOffer:
			f.assignManageBuyOfferResult(r.Tr.MustManageBuyOfferResult())
		case xdr.OperationTypeCreatePassiveSellOffer:
			f.assignManageSellOfferResult(r.Tr.MustCreatePassiveSellOfferResult())
		case xdr.OperationTypeSetOptions:
			f.assignSetOptionsResult(r.Tr.MustSetOptionsResult())
		case xdr.OperationTypeChangeTrust:
			f.assignChangeTrustResult(r.Tr.MustChangeTrustResult())
		case xdr.OperationTypeAllowTrust:
			f.assignAllowTrustResult(r.Tr.MustAllowTrustResult())
		case xdr.OperationTypeAccountMerge:
			f.assignAccountMergeResult(r.Tr.MustAccountMergeResult())
		case xdr.OperationTypeManageData:
			f.assignManageDataResult(r.Tr.MustManageDataResult())
		case xdr.OperationTypeBumpSequence:
			f.assignBumpSequenceResult(r.Tr.MustBumpSeqResult())
		case xdr.OperationTypeInflation:
			f.assignInflationResult(r.Tr.MustInflationResult())
		}
	} else {
		f.operation.Succesful = false
	}
}

func (f *operationFactory) assignCreateAccountResult(r xdr.CreateAccountResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.CreateAccountResultCodeCreateAccountSuccess
}

func (f *operationFactory) assignPaymentResult(r xdr.PaymentResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.PaymentResultCodePaymentSuccess
}

func (f *operationFactory) assignPathPaymentStrictReceiveResult(r xdr.PathPaymentStrictReceiveResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess

	if s, ok := r.GetSuccess(); ok {
		if len(s.Offers) > 0 {
			f.operation.AmountSent = amount.String(s.Offers[0].AmountBought)
		}
		f.operation.ResultLastAmount = amount.String(s.Last.Amount)
		f.operation.AmountReceived = f.operation.ResultLastAmount

		f.operation.ResultLastAsset = NewAsset(&s.Last.Asset)
		f.operation.ResultLastDestination = s.Last.Destination.Address()
	}

	if a, ok := r.GetNoIssuer(); ok {
		f.operation.ResultNoIssuer = NewAsset(&a)
	}
}

func (f *operationFactory) assignPathPaymentStrictSendResult(r xdr.PathPaymentStrictSendResult) {

	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess

	if s, ok := r.GetSuccess(); ok {
		if len(s.Offers) > 0 {
			f.operation.AmountSent = amount.String(s.Offers[0].AmountBought)
		}
		f.operation.ResultLastAmount = amount.String(s.Last.Amount)
		f.operation.AmountReceived = f.operation.ResultLastAmount

		f.operation.ResultLastAsset = NewAsset(&s.Last.Asset)
		f.operation.ResultLastDestination = s.Last.Destination.Address()
	}
	if a, ok := r.GetNoIssuer(); ok {
		f.operation.ResultNoIssuer = NewAsset(&a)
	}
}

func (f *operationFactory) assignManageSellOfferResult(r xdr.ManageSellOfferResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.ManageSellOfferResultCodeManageSellOfferSuccess

	if s, ok := r.GetSuccess(); ok {
		f.assignManageOfferResult(s)
	}
}

func (f *operationFactory) assignManageBuyOfferResult(r xdr.ManageBuyOfferResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.ManageBuyOfferResultCodeManageBuyOfferSuccess

	if s, ok := r.GetSuccess(); ok {
		f.assignManageOfferResult(s)
	}
}

func (f *operationFactory) assignManageOfferResult(s xdr.ManageOfferSuccessResult) {
	if o, ok := s.Offer.GetOffer(); ok {
		p, _ := big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()

		f.operation.ResultOffer = &Offer{
			Amount:   amount.String(o.Amount),
			Price:    p,
			PriceND:  Price{int(o.Price.N), int(o.Price.D)},
			Buying:   *NewAsset(&o.Buying),
			Selling:  *NewAsset(&o.Selling),
			OfferID:  int64(o.OfferId),
			SellerID: o.SellerId.Address(),
		}

		f.operation.ResultOfferEffect = strings.Replace(
			s.Offer.Effect.String(), "ManageOfferEffectManageOffer", "", 1,
		)
	}
}

func (f *operationFactory) assignSetOptionsResult(r xdr.SetOptionsResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.SetOptionsResultCodeSetOptionsSuccess
}

func (f *operationFactory) assignChangeTrustResult(r xdr.ChangeTrustResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.ChangeTrustResultCodeChangeTrustSuccess
}

func (f *operationFactory) assignAllowTrustResult(r xdr.AllowTrustResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.AllowTrustResultCodeAllowTrustSuccess
}

func (f *operationFactory) assignAccountMergeResult(r xdr.AccountMergeResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.AccountMergeResultCodeAccountMergeSuccess

	if b, ok := r.GetSourceAccountBalance(); ok {
		f.operation.ResultSourceAccountBalance = amount.String(b)
	}
}

func (f *operationFactory) assignManageDataResult(r xdr.ManageDataResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.ManageDataResultCodeManageDataSuccess
}

func (f *operationFactory) assignBumpSequenceResult(r xdr.BumpSequenceResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.BumpSequenceResultCodeBumpSequenceSuccess
}

func (f *operationFactory) assignInflationResult(r xdr.InflationResult) {
	f.operation.InnerResultCode = int(r.Code)
	f.operation.Succesful = r.Code == xdr.InflationResultCodeInflationSuccess
}
