package es

import (
	// "log"
	"time"

	"github.com/stellar/go/xdr"
)

// EventsExtractor is temporary struct holding data essential for processing the set of changes
type EventsExtractor struct {
	changes         []xdr.OperationMeta
	seq             int
	closeTime       time.Time
	basePagingToken PagingToken

	offerEvents []*OfferEvent
	index       int
}

// ProduceBalances constructs balance extracotr and returns balances
func ProduceEvents(changes []xdr.OperationMeta, ledgerSeq int, t time.Time, basePagingToken PagingToken) (offerEvents []*OfferEvent) {
	e := &EventsExtractor{
		changes:         changes,
		seq:             ledgerSeq,
		closeTime:       t,
		basePagingToken: basePagingToken,
		index:           0,
	}

	if e == nil {
		return offerEvents
	}

	return e.extract()
}

// Extract balances from current changes list
func (e *EventsExtractor) extract() []*OfferEvent {
	for _, meta := range e.changes {
		for _, change := range meta.Changes {
			switch t := change.Type; t {
			case xdr.LedgerEntryChangeTypeLedgerEntryState:
				e.state(change)

			case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
				e.created(change)

			case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
				e.updated(change)
			}
		}
	}

	return e.offerEvents
}

func (e *EventsExtractor) state(change xdr.LedgerEntryChange) {
	return

	// state := change.MustState().Data

	// switch x := state.Type; x {
	// case xdr.LedgerEntryTypeAccount:
	// 	account := state.MustAccount()
	// 	address := account.AccountId.Address()

	// 	e.values[address] = account.Balance
	// case xdr.LedgerEntryTypeTrustline:
	// 	line := state.MustTrustLine()
	// 	address := line.AccountId.Address()
	// 	e.values[address] = line.Balance
	// }
}

func (e *EventsExtractor) created(change xdr.LedgerEntryChange) {
	created := change.MustCreated().Data
	e.index++
	pagingToken := PagingToken{EffectIndex: e.index}.Merge(e.basePagingToken)

	switch x := created.Type; x {
	case xdr.LedgerEntryTypeOffer:
		offer := created.MustOffer()
		e.offerEvents = append(
			e.offerEvents,
			NewEventFromOfferEntry(EventTypeCreate, offer, e.seq, e.closeTime, pagingToken),
		)
	}
}

func (e *EventsExtractor) updated(change xdr.LedgerEntryChange) {
	updated := change.MustUpdated().Data

	e.index++
	pagingToken := PagingToken{EffectIndex: e.index}.Merge(e.basePagingToken)

	switch x := updated.Type; x {
	case xdr.LedgerEntryTypeOffer:
		offer := updated.MustOffer()

		e.offerEvents = append(
			e.offerEvents,
			NewEventFromOfferEntry(EventTypeUpdate, offer, e.seq, e.closeTime, pagingToken),
		)
	}
}
