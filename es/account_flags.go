package es

import "github.com/stellar/go/xdr"

// AccountFlags represents account flags from XDR
type AccountFlags struct {
	AuthRequired  bool `json:"required,omitempty"`
	AuthRevocable bool `json:"revocable,omitempty"`
	AuthImmutable bool `json:"immutable,omitempty"`
}

// NewAccountFlags returns account flags or nil if source flags are nil
func NewAccountFlags(flags *xdr.Uint32) *AccountFlags {
	if flags == nil {
		return nil
	}

	l := xdr.AccountFlags(*flags)

	return &AccountFlags{
		l&xdr.AccountFlagsAuthRequiredFlag != 0,
		l&xdr.AccountFlagsAuthRevocableFlag != 0,
		l&xdr.AccountFlagsAuthImmutableFlag != 0,
	}
}
