package es

import (
	"time"

	"github.com/stellar/go/xdr"
)

// AccountBalanceMap is account id <=> balance map
type AccountBalanceMap map[string]xdr.Int64

// BalanceExtractor is temporary struct holding data essential for processing the set of changes
type BalanceExtractor struct {
	Changes []xdr.LedgerEntryChange
	Time    time.Time
	Source  BalanceSource
	ID      string

	values   AccountBalanceMap
	balances []Balance
}

// NewBalanceExtractor constructs and returns instance of BalanceExtractor
func NewBalanceExtractor(changes []xdr.LedgerEntryChange, t time.Time, source BalanceSource, id string) *BalanceExtractor {
	return &BalanceExtractor{
		Changes: changes,
		Time:    t,
		Source:  source,
		ID:      id,
	}
}

// Extract balances from current changes list
func (e *BalanceExtractor) Extract() []Balance {
	for _, change := range e.Changes {
		switch t := change.Type; t {
		case xdr.LedgerEntryChangeTypeLedgerEntryState:
			e.state(change)

		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
			e.created(change)


		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
			e.updated(change)
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

func (e *BalanceExtractor) created(change xdr.LedgerEntryChange) {
	created := change.MustCreated().Data

	switch x := created.Type; x {
	case xdr.LedgerEntryTypeAccount:
		account := created.MustAccount()
		id := e.ID + ":" + account.AccountId.Address()

		e.balances = append(
			e.balances,
			*NewBalanceFromAccountEntry(account, e.Time, id, e.Source),
		)
	case xdr.LedgerEntryTypeTrustline:
		line := created.MustTrustLine()
		id := e.ID + ":" + line.AccountId.Address()

		e.balances = append(
			e.balances,
			*NewBalanceFromTrustLineEntry(line, e.Time, id, e.Source),
		)
	}
}

func (e *BalanceExtractor) updated(change xdr.LedgerEntryChange) {
	updated := change.MustUpdated().Data

	switch x := updated.Type; x {
	case xdr.LedgerEntryTypeAccount:
		account := updated.MustAccount()
		address := account.AccountId.Address()
		oldBalance := prev[address]

		if oldBalance != account.Balance {
			e.balances = append(
				e.balances, 
				*NewBalanceFromAccountEntry(account, e.Time, id, e.Source)
			)
		}
	case xdr.LedgerEntryTypeTrustline:
		line := updated.MustTrustLine()
		address := line.AccountId.Address()
		oldBalance := prev[address]

		if oldBalance != line.Balance {
			e.balances = append(
				balances,
				*NewBalanceFromTrustLineEntry(line, e.Time, id, e.Source)
			)
		}
	}
}

// // ExtractBalances returns balances extracted from metas
// func ExtractBalances(c []xdr.LedgerEntryChange, now time.Time, id string, source BalanceSource) []*Balance {
// 	var prev AccountBalanceMap
// 	var balances []*Balance

// 	for _, change := range c {
// 		switch t := change.Type; t {
// 		case xdr.LedgerEntryChangeTypeLedgerEntryState:
// 			processState(change, &prev)

// 		case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
// 			balances = processCreated(change, balances, now, id, source)

// 		case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
// 			updated := change.MustUpdated().Data

// 			switch x := updated.Type; x {
// 			case xdr.LedgerEntryTypeAccount:
// 				account := updated.MustAccount()
// 				oldBalance := prev[account.AccountId.Address()]
// 				if oldBalance != account.Balance {
// 					balances = append(balances, NewBalanceFromAccountEntry(account, now, id, source))
// 				}
// 			case xdr.LedgerEntryTypeTrustline:
// 				line := updated.MustTrustLine()
// 				oldBalance := prev[line.AccountId.Address()]
// 				if oldBalance != line.Balance {
// 					balances = append(balances, NewBalanceFromTrustLineEntry(line, now, id, source))
// 				}
// 			}
// 		}
// 	}

// 	return balances
// }

// func processState(change xdr.LedgerEntryChange, prev *AccountBalanceMap) {
// 	state := change.MustState().Data

// 	switch x := state.Type; x {
// 	case xdr.LedgerEntryTypeAccount:
// 		account := state.MustAccount()
// 		(*prev)[account.AccountId.Address()] = account.Balance
// 	case xdr.LedgerEntryTypeTrustline:
// 		line := state.MustTrustLine()
// 		(*prev)[line.AccountId.Address()] = line.Balance
// 	}
// }

// func processCreated(change xdr.LedgerEntryChange, balances []*Balance, now time.Time, id string, source BalanceSource) []*Balance {
// 	created := change.MustCreated().Data

// 	switch x := created.Type; x {
// 	case xdr.LedgerEntryTypeAccount:
// 		balances = append(balances, NewBalanceFromAccountEntry(created.MustAccount(), now, id, source))
// 	case xdr.LedgerEntryTypeTrustline:
// 		balances = append(balances, NewBalanceFromTrustLineEntry(created.MustTrustLine(), now, id, source))
// 	}
// 	return balances
// }
