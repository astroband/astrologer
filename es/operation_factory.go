package es

import (
	"math/big"
	"strings"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/xdr"
)

// OperationFactory represents operation builder
type operationFactory struct {
	transaction *Transaction
	source      *xdr.Operation
	result      *xdr.OperationResult
	index       int
	pagingToken PagingToken

	operation *Operation
}

// ProduceOperation creates factory and returns produced operation
func ProduceOperation(t *Transaction, o *xdr.Operation, r *xdr.OperationResult, n int) *Operation {
	factory := operationFactory{
		transaction: t,
		source:      o,
		result:      r,
		index:       n,
	}

	return factory.produce()
}

func (f *operationFactory) produce() *Operation {
	f.makeOperation()
	f.assignSourceAccountID()
	f.assignPagingToken()
	f.assignType()
	f.assignSpecifics()
	f.assignResult() // see operation_factory_result.go

	return f.operation
}

func (f *operationFactory) makeOperation() {
	f.operation = &Operation{
		TxID:              f.transaction.ID,
		TxIndex:           f.transaction.Index,
		Index:             f.index,
		Seq:               f.transaction.Seq,
		PagingToken:       f.pagingToken,
		CloseTime:         f.transaction.CloseTime,
		TxSourceAccountID: f.transaction.SourceAccountID,

		Memo: f.transaction.Memo,
	}
}

func (f *operationFactory) assignPagingToken() {
	f.operation.PagingToken = PagingToken{
		LedgerSeq:        f.transaction.Seq,
		TransactionOrder: f.transaction.Index,
		OperationOrder:   f.index,
	}
}

func (f *operationFactory) assignSourceAccountID() {
	sourceAccountID := f.transaction.SourceAccountID

	if f.source.SourceAccount != nil {
		sourceAccountID = f.source.SourceAccount.Address()
	}

	f.operation.SourceAccountID = sourceAccountID
}

func (f *operationFactory) assignType() {
	f.operation.Type = strings.Replace(f.source.Body.Type.String(), "OperationType", "", 1)
}

func (f *operationFactory) assignSpecifics() {
	body := f.source.Body

	switch t := body.Type; t {
	case xdr.OperationTypeCreateAccount:
		f.assignCreateAccount(body.MustCreateAccountOp())
	case xdr.OperationTypePayment:
		f.assignPayment(body.MustPaymentOp())
	case xdr.OperationTypePathPayment:
		f.assignPathPayment(body.MustPathPaymentOp())
	case xdr.OperationTypeManageSellOffer:
		f.assignManageSellOffer(body.MustManageSellOfferOp())
	case xdr.OperationTypeManageBuyOffer:
		f.assignManageBuyOffer(body.MustManageBuyOfferOp())
	case xdr.OperationTypeCreatePassiveSellOffer:
		f.assignCreatePassiveSellOffer(body.MustCreatePassiveSellOfferOp())
	case xdr.OperationTypeSetOptions:
		f.assignSetOptions(body.MustSetOptionsOp())
	case xdr.OperationTypeChangeTrust:
		f.assignChangeTrust(body.MustChangeTrustOp())
	case xdr.OperationTypeAllowTrust:
		f.assignAllowTrust(body.MustAllowTrustOp())
	case xdr.OperationTypeAccountMerge:
		f.assignAccountMerge(body.MustDestination())
	case xdr.OperationTypeManageData:
		f.assignManageData(body.MustManageDataOp())
	case xdr.OperationTypeBumpSequence:
		f.assignBumpSequence(body.MustBumpSequenceOp())
	}
}

func (f *operationFactory) assignCreateAccount(o xdr.CreateAccountOp) {
	f.operation.SourceAmount = amount.String(o.StartingBalance)
	f.operation.DestinationAccountID = o.Destination.Address()
}

func (f *operationFactory) assignPayment(o xdr.PaymentOp) {
	f.operation.SourceAmount = amount.String(o.Amount)
	f.operation.DestinationAccountID = o.Destination.Address()
	f.operation.SourceAsset = NewAsset(&o.Asset)
}

func (f *operationFactory) assignPathPayment(o xdr.PathPaymentOp) {
	f.operation.DestinationAccountID = o.Destination.Address()
	f.operation.DestinationAmount = amount.String(o.DestAmount)
	f.operation.DestinationAsset = NewAsset(&o.DestAsset)

	f.operation.SourceAmount = amount.String(o.SendMax)
	f.operation.SourceAsset = NewAsset(&o.SendAsset)

	f.operation.Path = make([]*Asset, len(o.Path))

	for i, a := range o.Path {
		f.operation.Path[i] = NewAsset(&a)
	}
}

func (f *operationFactory) assignManageSellOffer(o xdr.ManageSellOfferOp) {
	f.operation.SourceAmount = amount.String(o.Amount)
	f.operation.SourceAsset = NewAsset(&o.Buying)
	f.operation.OfferID = int(o.OfferId)
	f.operation.OfferPrice, _ = big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()
	f.operation.OfferPriceND = &Price{int(o.Price.N), int(o.Price.D)}
	f.operation.DestinationAsset = NewAsset(&o.Selling)
}

func (f *operationFactory) assignManageBuyOffer(o xdr.ManageBuyOfferOp) {
	f.operation.SourceAmount = amount.String(o.BuyAmount)
	f.operation.SourceAsset = NewAsset(&o.Selling)
	f.operation.DestinationAsset = NewAsset(&o.Buying)
	f.operation.OfferID = int(o.OfferId)
	f.operation.OfferPrice, _ = big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()
	f.operation.OfferPriceND = &Price{int(o.Price.N), int(o.Price.D)}
}

func (f *operationFactory) assignCreatePassiveSellOffer(o xdr.CreatePassiveSellOfferOp) {
	f.operation.SourceAmount = amount.String(o.Amount)
	f.operation.SourceAsset = NewAsset(&o.Buying)
	f.operation.OfferPrice, _ = big.NewRat(int64(o.Price.N), int64(o.Price.D)).Float64()
	f.operation.OfferPriceND = &Price{int(o.Price.N), int(o.Price.D)}
	f.operation.DestinationAsset = NewAsset(&o.Selling)
}

func (f *operationFactory) assignSetOptions(o xdr.SetOptionsOp) {
	if o.InflationDest != nil {
		f.operation.InflationDest = o.InflationDest.Address()
	}

	if o.HomeDomain != nil {
		f.operation.HomeDomain = string(*o.HomeDomain)
	}

	f.operation.Thresholds = NewAccountThresholds(
		o.LowThreshold, o.MedThreshold, o.HighThreshold, o.MasterWeight,
	)
	f.operation.SetFlags = NewAccountFlags(o.SetFlags)
	f.operation.ClearFlags = NewAccountFlags(o.ClearFlags)
	f.operation.Signer = NewSigner(o.Signer)
}

func (f *operationFactory) assignChangeTrust(o xdr.ChangeTrustOp) {
	f.operation.DestinationAmount = amount.String(o.Limit)
	f.operation.DestinationAsset = NewAsset(&o.Line)
}

func (f *operationFactory) assignAllowTrust(o xdr.AllowTrustOp) {
	a := o.Asset.ToAsset(o.Trustor)

	f.operation.DestinationAsset = NewAsset(&a)
	f.operation.DestinationAccountID = o.Trustor.Address()
	f.operation.Authorize = o.Authorize
}

func (f *operationFactory) assignAccountMerge(d xdr.AccountId) {
	f.operation.DestinationAccountID = d.Address()
}

func (f *operationFactory) assignBumpSequence(o xdr.BumpSequenceOp) {
	f.operation.BumpTo = int(o.BumpTo)
}

// TODO: Apply some magic to the value
func (f *operationFactory) assignManageData(o xdr.ManageDataOp) {
	f.operation.Data = &DataEntry{Name: string(o.DataName)}
	if o.DataValue != nil {
		f.operation.Data.Value = string(*o.DataValue)
	}
}
