package es

import (
	"encoding/json"
	"strconv"
)

// PagingToken represents numerical order / id of objects.
// Transaction 0 of the ledger 1 and the ledger itself will have the same order, so, start orders from 1.
type PagingToken struct {
	LedgerSeq        int
	TransactionOrder uint8
	OperationOrder   uint8
	AuxOrder1        uint8
	AuxOrder2        uint8
}

var (
	ledgerShift      uint = 32
	transactionShift uint = 24
	operationShift   uint = 16
	aux1Shift        uint = 8
	aux2Shift        uint
)

// UInt64 returns integers value from order
func (o PagingToken) UInt64() (result uint64) {
	result = result | (uint64(o.LedgerSeq) << ledgerShift)
	result = result | (uint64(o.TransactionOrder) << transactionShift)
	result = result | (uint64(o.OperationOrder) << operationShift)
	result = result | (uint64(o.AuxOrder1) << aux1Shift)
	result = result | (uint64(o.AuxOrder2) << aux2Shift)

	return result
}

// String returns string representation of order
func (o PagingToken) String() (result string) {
	return strconv.FormatUint(o.UInt64(), 10)
}

// MarshalJSON marshals to int
func (o PagingToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.UInt64())
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

	if o.AuxOrder1 != 0 {
		result.AuxOrder1 = o.AuxOrder1
	} else {
		result.AuxOrder1 = n.AuxOrder1
	}

	if o.AuxOrder2 != 0 {
		result.AuxOrder2 = o.AuxOrder2
	} else {
		result.AuxOrder2 = n.AuxOrder2
	}

	return result
}