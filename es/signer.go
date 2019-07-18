package es

import "github.com/stellar/go/xdr"

// Signer represents signer as export
type Signer struct {
	ID     string `json:"id"`
	Weight int    `json:"weight"`
	Type   int    `json:"type"`
}

// NewSigner returns new Signer
func NewSigner(signer *xdr.Signer) *Signer {
	if signer == nil {
		return nil
	}

	return &Signer{
		signer.Key.Address(),
		int(signer.Weight),
		int(signer.Key.Type),
	}
}
