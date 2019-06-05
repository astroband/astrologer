package es

import (
	"time"
)

// Trade represents trade entry
type Trade struct {
	ID              string    `json:"-"`
	Sold            string    `json:"sold"`
	Bought          string    `json:"bought"`
	AssetSold       Asset     `json:"asset_sold"`
	AssetBought     Asset     `json:"asset_bought"`
	SoldOfferID     int       `json:"sold_offer_id"`
	BoughtOfferID   int       `json:"bought_offer_id"`
	SellerID        string    `json:"seller_id"`
	BuyerID         string    `json:"buyer_id"`
	Price           string    `json:"price"`
	LedgerCloseTime time.Time `json:"ledger_close_time"`
}

// DocID balance es document id
func (t *Trade) DocID() *string {
	return &t.ID
}

// IndexName balances index name
func (t *Trade) IndexName() string {
	return tradesIndexName
}
