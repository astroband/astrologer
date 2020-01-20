package es

import (
	"encoding/json"
	"fmt"
)

// PagingToken represents numerical order / id of objects.
// Transaction 0 of the ledger 1 and the ledger itself will have the same order, so, start orders from 1.
type PagingToken struct {
	LedgerSeq        int
	TransactionOrder int
	OperationOrder   int
	EffectIndex      int
}

var (
	ledgerFormat      = "%012d"
	transactionFormat = "%04d"
	operationFormat   = "%04d"
	effectIndexFormat = "%04d"
)

// String returns string representation of order
func (o PagingToken) String() (result string) {
	return fmt.Sprintf(ledgerFormat, o.LedgerSeq) + "-" +
		fmt.Sprintf(transactionFormat, o.TransactionOrder) + "-" +
		fmt.Sprintf(operationFormat, o.OperationOrder) + "-" +
		fmt.Sprintf(effectIndexFormat, o.EffectIndex)
}

// MarshalJSON marshals to int
func (o PagingToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

// Merge merges with other order
func (o PagingToken) Merge(n PagingToken) (result PagingToken) {
	if o.LedgerSeq != 0 {
		result.LedgerSeq = o.LedgerSeq
	} else {
		result.LedgerSeq = n.LedgerSeq
	}

	if o.TransactionOrder != 0 {
		result.TransactionOrder = o.TransactionOrder
	} else {
		result.TransactionOrder = n.TransactionOrder
	}

	if o.OperationOrder != 0 {
		result.OperationOrder = o.OperationOrder
	} else {
		result.OperationOrder = n.OperationOrder
	}

	if o.EffectIndex != 0 {
		result.EffectIndex = o.EffectIndex
	} else {
		result.EffectIndex = n.EffectIndex
	}

	return result
}
