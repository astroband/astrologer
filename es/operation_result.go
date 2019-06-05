package es

import (
	"fmt"
	"math/big"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// AppendResult Appends operation result
func AppendResult(op *Operation, r *xdr.OperationResult) {
	op.ResultCode = int(r.Code)

	if r.Code == xdr.OperationResultCodeOpInner {
		switch t := r.Tr.Type; t {
		case xdr.OperationTypeCreateAccount:
			newCreateAccountResult(r.Tr.MustCreateAccountResult(), op)
		case xdr.OperationTypePayment:
			newPaymentResult(r.Tr.MustPaymentResult(), op)
		case xdr.OperationTypePathPayment:
			newPathPaymentResult(r.Tr.MustPathPaymentResult(), op)
		case xdr.OperationTypeManageSellOffer:
			newManageSellOfferResult(r.Tr.MustManageSellOfferResult(), op)
		case xdr.OperationTypeManageBuyOffer:
			newManageBuyOfferResult(r.Tr.MustManageBuyOfferResult(), op)
		case xdr.OperationTypeCreatePassiveSellOffer:
			newManageSellOfferResult(r.Tr.MustCreatePassiveSellOfferResult(), op)
		case xdr.OperationTypeSetOptions:
			newSetOptionsResult(r.Tr.MustSetOptionsResult(), op)
		case xdr.OperationTypeChangeTrust:
			newChangeTrustResult(r.Tr.MustChangeTrustResult(), op)
		case xdr.OperationTypeAllowTrust:
			newAllowTrustResult(r.Tr.MustAllowTrustResult(), op)
		case xdr.OperationTypeAccountMerge:
			newAccountMergeResult(r.Tr.MustAccountMergeResult(), op)
		case xdr.OperationTypeManageData:
			newManageDataResult(r.Tr.MustManageDataResult(), op)
		case xdr.OperationTypeBumpSequence:
			newBumpSequenceResult(r.Tr.MustBumpSeqResult(), op)
		}
	} else {
		op.Succesful = false
	}
}

func newCreateAccountResult(r xdr.CreateAccountResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.CreateAccountResultCodeCreateAccountSuccess
}

func newPaymentResult(r xdr.PaymentResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.PaymentResultCodePaymentSuccess
}

func newPathPaymentResult(r xdr.PathPaymentResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.PathPaymentResultCodePathPaymentSuccess

	if s, ok := r.GetSuccess(); ok {
		op.ResultOffersClaimed = appendOffersClaimed(s.Offers)

		if op.ResultOffersClaimed != nil {
			op.AmountSent = (*op.ResultOffersClaimed)[0].AmountBought
		}
		op.ResultLastAmount = amount.String(s.Last.Amount)
		op.AmountReceived = op.ResultLastAmount

		op.ResultLastAsset = NewAsset(&s.Last.Asset)
		op.ResultLastDestination = s.Last.Destination.Address()
	}

	if a, ok := r.GetNoIssuer(); ok {
		op.ResultNoIssuer = NewAsset(&a)
	}
}

func newManageSellOfferResult(r xdr.ManageSellOfferResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.ManageSellOfferResultCodeManageSellOfferSuccess

	if s, ok := r.GetSuccess(); ok {
		newManageOfferResult(s, op)
	}
}

func newManageBuyOfferResult(r xdr.ManageBuyOfferResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.ManageBuyOfferResultCodeManageBuyOfferSuccess

	if s, ok := r.GetSuccess(); ok {
		newManageOfferResult(s, op)
	}
}

func newManageOfferResult(s xdr.ManageOfferSuccessResult, op *Operation) {
	op.ResultOffersClaimed = appendOffersClaimed(s.OffersClaimed)

	if o, ok := s.Offer.GetOffer(); ok {
		p, _ := big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()

		op.ResultOffer = &Offer{
			Amount:   amount.String(o.Amount),
			Price:    p,
			PriceND:  Price{int(o.Price.N), int(o.Price.D)},
			Buying:   *NewAsset(&o.Buying),
			Selling:  *NewAsset(&o.Selling),
			OfferID:  int64(o.OfferId),
			SellerID: o.SellerId.Address(),
		}

		op.ResultOfferEffect = s.Offer.Effect.String()
	}
}

func newSetOptionsResult(r xdr.SetOptionsResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.SetOptionsResultCodeSetOptionsSuccess
}

func newChangeTrustResult(r xdr.ChangeTrustResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.ChangeTrustResultCodeChangeTrustSuccess
}

func newAllowTrustResult(r xdr.AllowTrustResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.AllowTrustResultCodeAllowTrustSuccess
}

func newAccountMergeResult(r xdr.AccountMergeResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.AccountMergeResultCodeAccountMergeSuccess

	if b, ok := r.GetSourceAccountBalance(); ok {
		op.ResultSourceAccountBalance = amount.String(b)
	}
}

func newManageDataResult(r xdr.ManageDataResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.ManageDataResultCodeManageDataSuccess
}

func newBumpSequenceResult(r xdr.BumpSequenceResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
	op.Succesful = r.Code == xdr.BumpSequenceResultCodeBumpSequenceSuccess
}

func appendOffersClaimed(c []xdr.ClaimOfferAtom) *[]OfferClaim {
	if len(c) > 0 {
		claims := make([]OfferClaim, len(c))
		for n := 0; n < len(c); n++ {
			c := c[n]

			fmt.Println("SOLD:", c.AssetSold, "BOUGHT:", c.AssetBought)

			claims[n] = OfferClaim{
				AmountSold:   amount.String(c.AmountSold),
				AmountBought: amount.String(c.AmountBought),
				AssetSold:    *NewAsset(&c.AssetSold),
				AssetBought:  *NewAsset(&c.AssetBought),
				OfferID:      int64(c.OfferId),
				SellerID:     c.SellerId.Address(),
			}
		}
		fmt.Println("---")

		return &claims
	}

	return nil
}
