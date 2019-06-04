package es

import (
	"time"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// BalanceSource represents type of balance record
type BalanceSource string

const (
	// BalanceSourceFee marks record came from tx_fee_meta
	BalanceSourceFee BalanceSource = "Fee"

	// BalanceSourceMeta marks record came from payment
	BalanceSourceMeta BalanceSource = "Meta"
)

// Balance represents balance log entry
type Balance struct {
	ID        string        `json:"-"`
	AccountID string        `json:"account_id"`
	Balance   string        `json:"balance"`
	Amount    string        `json:"amount"`
	CreatedAt time.Time     `json:"created_at"`
	Source    BalanceSource `json:"source"`
	Asset     Asset         `json:"asset"`
}

// NewBalanceFromAccountEntry creates Balance from AccountEntry
func NewBalanceFromAccountEntry(a xdr.AccountEntry, amt xdr.Int64, now time.Time, id string, source BalanceSource) *Balance {
	return &Balance{
		ID:        id,
		AccountID: a.AccountId.Address(),
		Balance:   amount.String(a.Balance),
		Amount:    amount.String(amt),
		Source:    source,
		CreatedAt: now,
		Asset:     *NewNativeAsset(),
	}
}

// NewBalanceFromTrustLineEntry creates Balance from TrustLineEntry
func NewBalanceFromTrustLineEntry(t xdr.TrustLineEntry, amt xdr.Int64, now time.Time, id string, source BalanceSource) *Balance {
	return &Balance{
		ID:        id,
		AccountID: t.AccountId.Address(),
		Balance:   amount.String(t.Balance),
		Amount:    amount.String(amt),
		Source:    source,
		CreatedAt: now,
		Asset:     *NewAsset(&t.Asset),
	}
}

// DocID balance es document id
func (b *Balance) DocID() *string {
	return &b.ID
}

// IndexName balances index name
func (b *Balance) IndexName() string {
	return balanceIndexName
}
