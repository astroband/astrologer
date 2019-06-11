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

	trades = e.fetchFromManageOffer(claims, e.operation.SourceAccountID, true)

	if e.operation.TxID == "278d693fde6e7ef879dc81779c858e479e96210d095d1ee3170ce35d432655ee" {
		fmt.Println(claims)
		for _, t := range trades {
			fmt.Println(t.AssetSold, t.AssetBought, t.Sold, t.Bought, t.Price, t.SellerID, t.BuyerID)
		}
	}

	return trades //e.fetchFromManageOffer(claims, e.operation.SourceAccountID, true)
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

	return e.fetchFromManageOffer(claims, e.operation.SourceAccountID, false)
}

func (e *TradeExtractor) fetchFromManageOffer(claims []xdr.ClaimOfferAtom, accountID string, sell bool) (trades []Trade) {
	trades = make([]Trade, len(claims))
	for i, claim := range claims {
		var price float64

		trade := Trade{
			PagingToken:     PagingToken{AuxOrder1: uint8(i)}.Merge(e.pagingToken),
			Sold:            amount.String(claim.AmountSold),
			Bought:          amount.String(claim.AmountBought),
			AssetSold:       *NewAsset(&claim.AssetSold),
			AssetBought:     *NewAsset(&claim.AssetBought),
			LedgerCloseTime: e.closeTime,
		}

		if sell {
			trade.SellerID = accountID
			trade.BuyerID = claim.SellerId.Address()
			price = float64(claim.AmountSold) / float64(claim.AmountBought)
		} else {
			trade.BuyerID = accountID
			trade.SellerID = claim.SellerId.Address()
			price = float64(claim.AmountBought) / float64(claim.AmountSold)
		}

		trade.Price = strconv.FormatFloat(price, 'f', 7, 64)

		trades[i] = trade
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
