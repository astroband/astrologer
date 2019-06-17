package es

import (
	"strconv"
	"time"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// TradeExtractor is used to extract trade index entries
type TradeExtractor struct {
	result      *xdr.OperationResult
	closeTime   time.Time
	pagingToken PagingToken
	operation   *Operation
	tokenIndex  int
}

// ProduceTrades returns trades
func ProduceTrades(r *xdr.OperationResult, op *Operation, closeTime time.Time, pagingToken PagingToken) (trades []Trade) {
	extractor := &TradeExtractor{
		result:      r,
		closeTime:   closeTime,
		pagingToken: pagingToken,
		operation:   op,
	}

	if extractor == nil {
		return trades
	}

	return extractor.extract()
}

// Extract returns trades entries
func (e *TradeExtractor) extract() (trades []Trade) {
	if e.result == nil {
		return trades
	}

	if e.result.Code == xdr.OperationResultCodeOpInner {
		switch t := e.result.Tr.Type; t {
		case xdr.OperationTypePathPayment:
			trades = e.fetchFromPathPayment(e.result.Tr.MustPathPaymentResult())
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
	if result.Code != xdr.ManageSellOfferResultCodeManageSellOfferSuccess {
		return trades
	}

	success, ok := result.GetSuccess()
	if !ok {
		return trades
	}

	claims := success.OffersClaimed
	if len(claims) == 0 {
		return trades
	}

	return e.fetchClaims(claims, e.operation.SourceAccountID)
}

func (e *TradeExtractor) fetchFromManageBuyOffer(result xdr.ManageBuyOfferResult) (trades []Trade) {
	if result.Code != xdr.ManageBuyOfferResultCodeManageBuyOfferSuccess {
		return trades
	}

	success, ok := result.GetSuccess()
	if !ok {
		return trades
	}

	claims := success.OffersClaimed
	if len(claims) == 0 {
		return trades
	}

	return e.fetchClaims(claims, e.operation.SourceAccountID)
}

func (e *TradeExtractor) fetchFromPathPayment(result xdr.PathPaymentResult) (trades []Trade) {
	if result.Code != xdr.PathPaymentResultCodePathPaymentSuccess {
		return trades
	}

	success, ok := result.GetSuccess()
	if !ok {
		return trades
	}

	claims := success.Offers
	if len(claims) == 0 {
		return trades
	}

	return e.fetchClaims(claims, e.operation.SourceAccountID)
}

func (e *TradeExtractor) fetchClaims(claims []xdr.ClaimOfferAtom, accountID string) (trades []Trade) {
	for _, claim := range claims {
		pagingTokenA := PagingToken{EffectGroup: TradeEffectPagingTokenGroup, EffectIndex: e.tokenIndex + 1}.Merge(e.pagingToken)

		tradeA := Trade{
			PagingToken:     pagingTokenA,
			OfferID:         int64(claim.OfferId),
			LedgerCloseTime: e.closeTime,
		}

		pagingTokenB := PagingToken{EffectGroup: TradeEffectPagingTokenGroup, EffectIndex: e.tokenIndex + 2}.Merge(e.pagingToken)

		tradeB := Trade{
			PagingToken:     pagingTokenB,
			OfferID:         int64(claim.OfferId),
			LedgerCloseTime: e.closeTime,
		}

		e.tokenIndex += 2

		tradeA.Sold = amount.String(claim.AmountSold)
		tradeA.Bought = amount.String(claim.AmountBought)
		tradeA.AssetSold = *NewAsset(&claim.AssetSold)
		tradeA.AssetBought = *NewAsset(&claim.AssetBought)
		tradeA.SellerID = accountID
		tradeA.BuyerID = claim.SellerId.Address()
		tradeA.Price = strconv.FormatFloat(float64(claim.AmountSold)/float64(claim.AmountBought), 'f', 7, 64)

		tradeB.Sold = amount.String(claim.AmountBought)
		tradeB.Bought = amount.String(claim.AmountSold)
		tradeB.AssetSold = *NewAsset(&claim.AssetBought)
		tradeB.AssetBought = *NewAsset(&claim.AssetSold)
		tradeB.SellerID = claim.SellerId.Address()
		tradeB.BuyerID = accountID
		tradeB.Price = strconv.FormatFloat(float64(claim.AmountBought)/float64(claim.AmountSold), 'f', 7, 64)

		trades = append(trades, tradeA)
		trades = append(trades, tradeB)
	}

	return trades
}
