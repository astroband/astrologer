package es

// Offer represents offer in ManageOffer
type Offer struct {
	Amount   string  `json:"amount"`
	Price    float64 `json:"price"`
	PriceND  Price   `json:"price_n_d"`
	Selling  Asset   `json:"selling"`
	Buying   Asset   `json:"buying"`
	OfferID  int64   `json:"offer_id"`
	SellerID string  `json:"seller_id"`
}
