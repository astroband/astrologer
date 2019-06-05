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
	Order     Order         `json:"order"`
	AccountID string        `json:"account_id"`
	Value     string        `json:"value"`
	Diff      string        `json:"diff"`
	CreatedAt time.Time     `json:"created_at"`
	Source    BalanceSource `json:"source"`
	Asset     Asset         `json:"asset"`
}

// NewBalanceFromAccountEntry creates Balance from AccountEntry
func NewBalanceFromAccountEntry(a xdr.AccountEntry, diff xdr.Int64, now time.Time, order Order, source BalanceSource) *Balance {
	return &Balance{
		Order:     order,
		AccountID: a.AccountId.Address(),
		Value:     amount.String(a.Balance),
		Diff:      amount.String(diff),
		Source:    source,
		CreatedAt: now,
		Asset:     *NewNativeAsset(),
	}
}

// NewBalanceFromTrustLineEntry creates Balance from TrustLineEntry
func NewBalanceFromTrustLineEntry(t xdr.TrustLineEntry, diff xdr.Int64, now time.Time, order Order, source BalanceSource) *Balance {
	return &Balance{
		Order:     order,
		AccountID: t.AccountId.Address(),
		Value:     amount.String(t.Balance),
		Diff:      amount.String(diff),
		Source:    source,
		CreatedAt: now,
		Asset:     *NewAsset(&t.Asset),
	}
}

// DocID balance es document id
func (b *Balance) DocID() *string {
	s := b.Order.String()
	return &s
}

// IndexName balances index name
func (b *Balance) IndexName() string {
	return balanceIndexName
}
