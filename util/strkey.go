package util

import (
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

func Address(a xdr.MuxedAccount) (string, error) {
	var (
		accountIdBin []byte
		marshalErr   error
	)

	if a.Type == xdr.CryptoKeyTypeKeyTypeEd25519 {
		accountIdBin, marshalErr = a.Ed25519.MarshalBinary()
	} else {
		accountIdBin, marshalErr = a.Med25519.Ed25519.MarshalBinary()
	}

	if marshalErr != nil {
		return "", marshalErr
	}

	return strkey.Encode(strkey.VersionByteAccountID, accountIdBin)
}
