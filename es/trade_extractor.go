package es

import (
	"fmt"
	"time"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

type TradeExtractor struct {
	result      *xdr.OperationResult
	index       int
	closeTime   time.Time
	pagingToken PagingToken
}

func NewTradeExtractor(r *[]xdr.OperationResult, index int, closeTime time.Time, pagingToken PagingToken) *TradeExtractor {
	if r == nil {
		return nil
	}

	result := (*r)[index]

	return &TradeExtractor{
		result:      &result,
		index:       index,
		closeTime:   closeTime,
		pagingToken: pagingToken,
	}
}

func (e *TradeExtractor) Extract() (trades []Trade) {
	if e.result == nil {
		return trades
	}

	if e.result.Code == xdr.OperationResultCodeOpInner {
		switch t := e.result.Tr.Type; t {
		// case xdr.OperationTypePathPayment:
		// 	result = newPathPaymentResult(r.Tr.MustPathPaymentResult())
		case xdr.OperationTypeManageSellOffer:
			trades = e.fetchFromManageSellOffer(e.result.Tr.MustManageSellOfferResult())
		case xdr.OperationTypeCreatePassiveSellOffer:
			trades = e.fetchFromManageSellOffer(e.result.Tr.MustCreatePassiveSellOfferResult())
		case xdr.OperationTypeManageBuyOffer:
			trades = e.fetchFromManageBuyOffer(e.result.Tr.MustManageBuyOfferResult())
		}
	}

	return trades
}

func (e *TradeExtractor) fetchFromManageSellOffer(result xdr.ManageSellOfferResult) (trades []Trade) {
	success, ok := result.GetSuccess()
	if !ok {
		return trades
	}

	offer, ok := success.Offer.GetOffer()
	if !ok {
		return trades
	}

	claims := success.OffersClaimed
	if len(claims) == 0 {
		return trades
	}

	return e.fetchFromOffer(offer, claims)
}

func (e *TradeExtractor) fetchFromManageBuyOffer(result xdr.ManageBuyOfferResult) (trades []Trade) {
	success, ok := result.GetSuccess()
	if !ok {
		return trades
	}

	offer, ok := success.Offer.GetOffer()
	if !ok {
		return trades
	}

	claims := success.OffersClaimed
	if len(claims) == 0 {
		return trades
	}

	return e.fetchFromOffer(offer, claims)
}

func (e *TradeExtractor) fetchFromOffer(offer xdr.OfferEntry, claims []xdr.ClaimOfferAtom) (trades []Trade) {
	trades = make([]Trade, len(claims))
	for i, claim := range claims {
		price := float64(claim.AmountSold) / float64(claim.AmountBought)

		trades[i] = Trade{
			PagingToken:     PagingToken{AuxOrder1: uint8(i)}.Merge(e.pagingToken),
			Sold:            amount.String(claim.AmountSold),
			Bought:          amount.String(claim.AmountBought),
			AssetSold:       *NewAsset(&claim.AssetSold),
			AssetBought:     *NewAsset(&claim.AssetBought),
			SoldOfferID:     int64(offer.OfferId),
			BoughtOfferID:   int64(claim.OfferId),
			SellerID:        offer.SellerId.Address(),
			BuyerID:         claim.SellerId.Address(),
			Price:           fmt.Sprintf("%f", price),
			LedgerCloseTime: e.closeTime,
		}
	}

	return trades
}
