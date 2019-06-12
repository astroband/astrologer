package es

import (
	"fmt"
	"strconv"
	"time"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// TradeExtractor is used to extract trade index entries
type TradeExtractor struct {
	result      *xdr.OperationResult
	index       int
	closeTime   time.Time
	pagingToken PagingToken
	operation   *Operation
}

// NewTradeExtractor creates new TradesExtractor or returns nil if result is inappropriate
func NewTradeExtractor(r *[]xdr.OperationResult, op *Operation, index int, closeTime time.Time, pagingToken PagingToken) *TradeExtractor {
	if r == nil {
		return nil
	}

	result := (*r)[index]

	return &TradeExtractor{
		result:      &result,
		index:       index,
		closeTime:   closeTime,
		pagingToken: pagingToken,
		operation:   op,
	}
}

// Extract returns trades entries
func (e *TradeExtractor) Extract() (trades []Trade) {
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

	return e.fetchFromManageOffer(claims, e.operation.SourceAccountID)
}

func (e *TradeExtractor) fetchFromManageBuyOffer(result xdr.ManageBuyOfferResult) (trades []Trade) {
	fmt.Println("BUUUUUUUUY")
	fmt.Println(e.operation.TxID)

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

	return e.fetchFromManageOffer(claims, e.operation.SourceAccountID)
}

func (e *TradeExtractor) fetchFromManageOffer(claims []xdr.ClaimOfferAtom, accountID string) (trades []Trade) {
	for i, claim := range claims {
		tradeA := Trade{
			PagingToken:     PagingToken{AuxOrder1: uint8(i)}.Merge(e.pagingToken),
			OfferID:         int64(claim.OfferId),
			LedgerCloseTime: e.closeTime,
		}

		tradeB := Trade{
			PagingToken:     PagingToken{AuxOrder1: uint8(i)}.Merge(e.pagingToken),
			OfferID:         int64(claim.OfferId),
			LedgerCloseTime: e.closeTime,
		}

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

func (e *TradeExtractor) fetchFromPathPayment(result xdr.PathPaymentResult) (trades []Trade) {
	return trades
	// if result.Code != xdr.PathPaymentResultCodePathPaymentSuccess {
	// 	return trades
	// }

	// success, ok := result.GetSuccess()
	// if !ok {
	// 	return trades
	// }

	// claims := success.Offers
	// if len(claims) == 0 {
	// 	return trades
	// }

	// for _, c := range claims {
	// 	fmt.Println(c.AssetSold, c.AssetBought)
	// }

	// return trades
}
