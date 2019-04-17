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
	AccountID string        `json:"account_id"`
	Balance   string        `json:"balance"`
	CreatedAt time.Time     `json:"created_at"`
	Source    BalanceSource `json:"source"`
	Asset     Asset         `json:"asset"`
}

// NewBalanceFromAccountEntry creates Balance from AccountEntry
func NewBalanceFromAccountEntry(a xdr.AccountEntry, now time.Time) *Balance {
	return &Balance{
		AccountID: a.AccountId.Address(),
		Balance:   amount.String(a.Balance),
		Source:    BalanceSourceMeta,
		CreatedAt: now,
		Asset:     *NewNativeAsset(),
	}
}

// NewBalanceFromTrustLineEntry creates Balance from TrustLineEntry
func NewBalanceFromTrustLineEntry(t xdr.TrustLineEntry, now time.Time) *Balance {
	return &Balance{
		AccountID: t.AccountId.Address(),
		Balance:   amount.String(t.Balance),
		Source:    BalanceSourceMeta,
		CreatedAt: now,
		Asset:     *NewAsset(&t.Asset),
	}
}

// ExtractBalances returns balances extracted from metas
func ExtractBalances(c []xdr.LedgerEntryChange, now time.Time) []*Balance {
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
				balances = append(balances, NewBalanceFromAccountEntry(created.MustAccount(), now))
			case xdr.LedgerEntryTypeTrustline:
				balances = append(balances, NewBalanceFromTrustLineEntry(created.MustTrustLine(), now))
			}

		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			updated := change.MustUpdated().Data

			switch x := updated.Type; x {
			case xdr.LedgerEntryTypeAccount:
				account := updated.MustAccount()
				oldBalance := prev[account.AccountId.Address()]
				if oldBalance != account.Balance {
					balances = append(balances, NewBalanceFromAccountEntry(account, now))
				}
			case xdr.LedgerEntryTypeTrustline:
				line := updated.MustTrustLine()
				oldBalance := prev[line.AccountId.Address()]
				if oldBalance != line.Balance {
					balances = append(balances, NewBalanceFromTrustLineEntry(line, now))
				}
			}
		}
	}

	return balances
}

// TODO: Balance key MUST exist (dups)
// DocID balance es document id
func (b *Balance) DocID() *string {
	return nil
}

// IndexName balances index name
func (b *Balance) IndexName() string {
	return balanceIndexName
}
