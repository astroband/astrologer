package es

import (
	"time"

	"github.com/stellar/go/xdr"
)

// AccountBalanceMap is account id <=> balance map
type AccountBalanceMap map[string]xdr.Int64

// BalanceExtractor is temporary struct holding data essential for processing the set of changes
type BalanceExtractor struct {
	changes         []xdr.LedgerEntryChange
	closeTime       time.Time
	source          BalanceSource
	basePagingToken PagingToken

	values   AccountBalanceMap
	balances []*Balance
}

// NewBalanceExtractor constructs and returns instance of BalanceExtractor
func NewBalanceExtractor(changes []xdr.LedgerEntryChange, t time.Time, source BalanceSource, basePagingToken PagingToken) *BalanceExtractor {
	return &BalanceExtractor{
		changes:         changes,
		closeTime:       t,
		source:          source,
		basePagingToken: basePagingToken,
		values:          make(AccountBalanceMap),
	}
}

// Extract balances from current changes list
func (e *BalanceExtractor) Extract() []*Balance {
	for n, change := range e.changes {
		switch t := change.Type; t {
		case xdr.LedgerEntryChangeTypeLedgerEntryState:
			e.state(change)

		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			e.created(change, byte(n+1))

		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			e.updated(change, byte(n+1))
		}
	}

	return e.balances
}

func (e *BalanceExtractor) state(change xdr.LedgerEntryChange) {
	state := change.MustState().Data

	switch x := state.Type; x {
	case xdr.LedgerEntryTypeAccount:
		account := state.MustAccount()
		address := account.AccountId.Address()

		e.values[address] = account.Balance
	case xdr.LedgerEntryTypeTrustline:
		line := state.MustTrustLine()
		address := line.AccountId.Address()
		e.values[address] = line.Balance
	}
}

func (e *BalanceExtractor) created(change xdr.LedgerEntryChange, n byte) {
	created := change.MustCreated().Data
	pagingToken := PagingToken{AuxOrder2: n}.Merge(e.basePagingToken)

	switch x := created.Type; x {
	case xdr.LedgerEntryTypeAccount:
		account := created.MustAccount()

		e.balances = append(
			e.balances,
			NewBalanceFromAccountEntry(account, account.Balance, e.closeTime, pagingToken, e.source),
		)
	case xdr.LedgerEntryTypeTrustline:
		line := created.MustTrustLine()

		e.balances = append(
			e.balances,
			NewBalanceFromTrustLineEntry(line, line.Balance, e.closeTime, pagingToken, e.source),
		)
	}
}

func (e *BalanceExtractor) updated(change xdr.LedgerEntryChange, n byte) {
	updated := change.MustUpdated().Data

	pagingToken := PagingToken{AuxOrder2: n}.Merge(e.basePagingToken)

	switch x := updated.Type; x {
	case xdr.LedgerEntryTypeAccount:
		account := updated.MustAccount()
		address := account.AccountId.Address()
		oldBalance := e.values[address]

		if oldBalance != account.Balance {
			diff := account.Balance - oldBalance

			e.balances = append(
				e.balances,
				NewBalanceFromAccountEntry(account, diff, e.closeTime, pagingToken, e.source),
			)
		}
	case xdr.LedgerEntryTypeTrustline:
		line := updated.MustTrustLine()
		address := line.AccountId.Address()
		oldBalance := e.values[address]

		if oldBalance != line.Balance {
			diff := line.Balance - oldBalance

			e.balances = append(
				e.balances,
				NewBalanceFromTrustLineEntry(line, diff, e.closeTime, pagingToken, e.source),
			)
		}
	}
}