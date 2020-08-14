package support

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/guregu/null"
	"github.com/stellar/go/xdr"
)

func MemoValue(memo xdr.Memo) null.String {
	var (
		value string
		valid bool
	)
	switch memo.Type {
	case xdr.MemoTypeMemoNone:
		value, valid = "", false
	case xdr.MemoTypeMemoText:
		scrubbed := Utf8Scrub(memo.MustText())
		notnull := strings.Join(strings.Split(scrubbed, "\x00"), "")
		value, valid = notnull, true
	case xdr.MemoTypeMemoId:
		value, valid = fmt.Sprintf("%d", memo.MustId()), true
	case xdr.MemoTypeMemoHash:
		hash := memo.MustHash()
		value, valid =
			base64.StdEncoding.EncodeToString(hash[:]),
			true
	case xdr.MemoTypeMemoReturn:
		hash := memo.MustRetHash()
		value, valid =
			base64.StdEncoding.EncodeToString(hash[:]),
			true
	default:
		panic(fmt.Errorf("invalid memo type: %v", memo.Type))
	}

	return null.NewString(value, valid)
}

func OperationMeta(txMeta xdr.TransactionMeta, opIndex int) *xdr.OperationMeta {
	if v1, ok := txMeta.GetV1(); ok {
		ops := v1.Operations
		return &ops[opIndex]
	}

	ops, ok := txMeta.GetOperations()
	if !ok {
		return nil
	}

	return &ops[opIndex]
}

// Copy paste from Horizon
func Utf8Scrub(in string) string {

	// First check validity using the stdlib, returning if the string is already
	// valid
	if utf8.ValidString(in) {
		return in
	}

	left := []byte(in)
	var result bytes.Buffer

	for len(left) > 0 {
		r, n := utf8.DecodeRune(left)

		_, err := result.WriteRune(r)
		if err != nil {
			panic(err)
		}

		left = left[n:]
	}

	return result.String()
}
