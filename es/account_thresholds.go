package es

import "github.com/stellar/go/xdr"

// AccountThresholds represents account thresholds from XDR
type AccountThresholds struct {
	Low    *byte `json:"low,omitempty"`
	Medium *byte `json:"medium,omitempty"`
	High   *byte `json:"high,omitempty"`
	Master *byte `json:"master,omitempty"`
}

// NewAccountThresholds constructs account thresholds from values, or returns nil if thresholds missing
func NewAccountThresholds(low *xdr.Uint32, medium *xdr.Uint32, high *xdr.Uint32, master *xdr.Uint32) *AccountThresholds {
	var thresholds *AccountThresholds

	if (low != nil) || (medium != nil) || (high != nil) || (master != nil) {
		thresholds = &AccountThresholds{}

		if low != nil {
			thresholds.Low = new(byte)
			*thresholds.Low = byte(*low)
		}

		if medium != nil {
			thresholds.Medium = new(byte)
			*thresholds.Medium = byte(*medium)
		}

		if high != nil {
			thresholds.High = new(byte)
			*thresholds.High = byte(*high)
		}

		if master != nil {
			thresholds.Master = new(byte)
			*thresholds.Master = byte(*master)
		}
	}

	return thresholds
}
