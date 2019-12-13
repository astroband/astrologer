package es

import (
	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
	"time"
)

type EventType string
type EventEntity string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeRemove EventType = "remove"
)

const (
	EventTypeOffer EventEntity = "offer"
)

// Event represents events log entry
type Event struct {
	Seq         int         `json:"seq"`
	PagingToken PagingToken `json:"paging_token"`
	AccountID   string      `json:"account_id"`
	CreatedAt   time.Time   `json:"created_at"`
	Entity      EventEntity `json:"entity"`
	Type        EventType   `json:"type"`
}

type OfferEvent struct {
	*Event
	*Offer
}

func NewEventFromOfferEntry(t EventType, offer xdr.OfferEntry, seq int, now time.Time, pagingToken PagingToken) *OfferEvent {
	price := Price{
		N: int(offer.Price.N),
		D: int(offer.Price.D),
	}
	return &OfferEvent{
		Event: &Event{
			Seq:         seq,
			PagingToken: pagingToken,
			AccountID:   offer.SellerId.Address(),
			CreatedAt:   now,
			Entity:      EventTypeOffer,
			Type:        t,
		},
		Offer: &Offer{
			Amount:   amount.String(offer.Amount),
			PriceND:  price,
			Price:    price.float(),
			Selling:  *NewAsset(&offer.Selling),
			Buying:   *NewAsset(&offer.Buying),
			OfferID:  int64(offer.OfferId),
			SellerID: offer.SellerId.Address(),
		},
	}
}

// DocID balance es document id
func (b *OfferEvent) DocID() *string {
	s := b.PagingToken.String()
	return &s
}

// IndexName balances index name
func (b *OfferEvent) IndexName() string {
	return eventIndexName
}
