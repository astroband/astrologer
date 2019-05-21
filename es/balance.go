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

	// BalanceSourceMeta marks record came from tx_meta
	BalanceSourceMeta BalanceSource = "Meta"
)

// Balance represents balance log entry
type Balance struct {
	ID        string        `json:"-"`
	AccountID string        `json:"account_id"`
	Balance   string        `json:"balance"`
	CreatedAt time.Time     `json:"created_at"`
	Source    BalanceSource `json:"source"`
	Asset     Asset         `json:"asset"`
}

// NewBalanceFromAccountEntry creates Balance from AccountEntry
func NewBalanceFromAccountEntry(a xdr.AccountEntry, now time.Time, id string, source BalanceSource) *Balance {
	return &Balance{
		ID:        id,
		AccountID: a.AccountId.Address(),
		Balance:   amount.String(a.Balance),
		Source:    source,
		CreatedAt: now,
		Asset:     *NewNativeAsset(),
	}
}

// NewBalanceFromTrustLineEntry creates Balance from TrustLineEntry
func NewBalanceFromTrustLineEntry(t xdr.TrustLineEntry, now time.Time, id string, source BalanceSource) *Balance {
	return &Balance{
		ID:        id,
		AccountID: t.AccountId.Address(),
		Balance:   amount.String(t.Balance),
		Source:    source,
		CreatedAt: now,
		Asset:     *NewAsset(&t.Asset),
	}
}

// ExtractBalances returns balances extracted from metas
func ExtractBalances(c []xdr.LedgerEntryChange, now time.Time, id string, source BalanceSource) []*Balance {
	var prev = make(map[string]xdr.Int64)
	var balances []*Balance

	for _, change := range c {
		switch t := change.Type; t {
		case xdr.LedgerEntryChangeTypeLedgerEntryState:
			state := change.MustState().Data

			switch x := state.Type; x {
			case xdr.LedgerEntryTypeAccount:
				account := state.MustAccount()
				prev[account.AccountId.Address()] = account.Balance
			case xdr.LedgerEntryTypeTrustline:
				line := state.MustTrustLine()
				prev[line.AccountId.Address()] = line.Balance
			}

		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			created := change.MustCreated().Data

			switch x := created.Type; x {
			case xdr.LedgerEntryTypeAccount:
				balances = append(balances, NewBalanceFromAccountEntry(created.MustAccount(), now, id, source))
			case xdr.LedgerEntryTypeTrustline:
				balances = append(balances, NewBalanceFromTrustLineEntry(created.MustTrustLine(), now, id, source))
			}

		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			updated := change.MustUpdated().Data

			switch x := updated.Type; x {
			case xdr.LedgerEntryTypeAccount:
				account := updated.MustAccount()
				oldBalance := prev[account.AccountId.Address()]
				if oldBalance != account.Balance {
					balances = append(balances, NewBalanceFromAccountEntry(account, now, id, source))
				}
			case xdr.LedgerEntryTypeTrustline:
				line := updated.MustTrustLine()
				oldBalance := prev[line.AccountId.Address()]
				if oldBalance != line.Balance {
					balances = append(balances, NewBalanceFromTrustLineEntry(line, now, id, source))
				}
			}
		}
	}

	return balances
}

// DocID balance es document id
func (b *Balance) DocID() *string {
	return &b.ID
}

// IndexName balances index name
func (b *Balance) IndexName() string {
	return balanceIndexName
}
