package es

import (
	"github.com/stellar/go/xdr"
)

// AppendResult Appends operation result
func AppendResult(op *Operation, r *xdr.OperationResult) {
	op.Succesful = r.Code == xdr.OperationResultCodeOpInner
	op.ResultCode = int(r.Code)

	if op.Succesful {
		switch t := r.Tr.Type; t {
		case xdr.OperationTypeCreateAccount:
			newCreateAccountResult(r.Tr.MustCreateAccountResult(), op)
		case xdr.OperationTypePayment:
			newPaymentResult(r.Tr.MustPaymentResult(), op)
			// case xdr.OperationTypePathPayment:
			// 	newPathPayment(o.Body.MustPathPaymentOp(), op)
			// case xdr.OperationTypeManageOffer:
			// 	newManageOffer(o.Body.MustManageOfferOp(), op)
			// case xdr.OperationTypeCreatePassiveOffer:
			// 	newCreatePassiveOffer(o.Body.MustCreatePassiveOfferOp(), op)
			// case xdr.OperationTypeSetOptions:
			// 	newSetOptions(o.Body.MustSetOptionsOp(), op)
			// case xdr.OperationTypeChangeTrust:
			// 	newChangeTrust(o.Body.MustChangeTrustOp(), op)
			// case xdr.OperationTypeAllowTrust:
			// 	newAllowTrust(o.Body.MustAllowTrustOp(), op)
			// case xdr.OperationTypeAccountMerge:
			// 	newAccountMerge(o.Body.MustDestination(), op)
			// case xdr.OperationTypeManageData:
			// 	newManageData(o.Body.MustManageDataOp(), op)
			// case xdr.OperationTypeBumpSequence:
			// 	newBumpSequence(o.Body.MustBumpSequenceOp(), op)
		}
	}
}

func newCreateAccountResult(r xdr.CreateAccountResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
}

func newPaymentResult(r xdr.PaymentResult, op *Operation) {
	op.InnerResultCode = int(r.Code)
}
