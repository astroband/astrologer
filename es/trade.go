package es

import (
	"time"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// Trade represents trade entry
type Trade struct {
	PagingToken     PagingToken `json:"paging_token"`
	Sold            string      `json:"sold"`
	Bought          string      `json:"bought"`
	AssetSold       Asset       `json:"asset_sold"`
	AssetBought     Asset       `json:"asset_bought"`
	SoldOfferID     int         `json:"sold_offer_id"`
	BoughtOfferID   int         `json:"bought_offer_id"`
	SellerID        string      `json:"seller_id"`
	BuyerID         string      `json:"buyer_id"`
	Price           string      `json:"price"`
	LedgerCloseTime time.Time   `json:"ledger_close_time"`
}

func NewTradesFromResult(r *xdr.OperationResult, opIndex int) (trades []Trade) {
	if r == nil {
		return trades
	}

	if r.Code == xdr.OperationResultCodeOpInner {
		switch t := r.Tr.Type; t {
		// case xdr.OperationTypePathPayment:
		// 	result = newPathPaymentResult(r.Tr.MustPathPaymentResult())
		case xdr.OperationTypeManageSellOffer:
			trades = fetchTradesFromManageSellOffer(r.Tr.MustManageSellOfferResult())
			// case xdr.OperationTypeManageBuyOffer:
			// 	result = fetchTradesFromManageBuyOffer(r.Tr.MustManageBuyOfferResult())
			// case xdr.OperationTypeCreatePassiveSellOffer:
			// 	result = newManageSellOfferResult(r.Tr.MustCreatePassiveSellOfferResult())
		}
	}

	return trades
}

func fetchTradesFromManageSellOffer(r xdr.ManageSellOfferResult) (trades []Trade) {
	success, ok := r.GetSuccess()
	if !ok {
		return trades
	}

	claims := success.OffersClaimed
	if len(claims) == 0 {
		return trades
	}

	trades = make([]Trade, len(claims))
	for i, claim := range claims {
		trades[i] = Trade{
			PagingToken: PagingToken{},
			Sold:        amount.String(claim.AmountSold),
			Bought:      amount.String(claim.AmountBought),
			AssetSold:   *NewAsset(&claim.AssetSold),
			AssetBought: *NewAsset(&claim.AssetBought),
			// SoldOfferID:
			// BoguthOfferID:
			// SellerID:
			// BuyerID:
			// Price:
			// LedgerCloseTime:
		}
	}

	return trades
}

// DocID balance es document id
func (t *Trade) DocID() *string {
	s := t.PagingToken.String()
	return &s
}

// IndexName balances index name
func (t *Trade) IndexName() string {
	return tradesIndexName
}
