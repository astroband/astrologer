package es

import (
	"math/big"
	"strings"

	"github.com/astroband/astrologer/util"
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
func ProduceOperation(t *Transaction, o *xdr.Operation, r *xdr.OperationResult, n int) (*Operation, error) {
	factory := operationFactory{
		transaction: t,
		source:      o,
		result:      r,
		index:       n,
	}

	return factory.produce()
}

func (f *operationFactory) produce() (*Operation, error) {
	var err error

	f.makeOperation()
	err = f.assignSourceAccountID()

	if err != nil {
		return nil, err
	}

	f.assignPagingToken()
	f.assignType()
	err = f.assignSpecifics()

	if err != nil {
		return nil, err
	}

	f.assignResult() // see operation_factory_result.go

	return f.operation, nil
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

func (f *operationFactory) assignSourceAccountID() error {
	var err error
	sourceAccountID := f.transaction.SourceAccountID

	if f.source.SourceAccount != nil {
		sourceAccountID, err = util.EncodeMuxedAccount(*f.source.SourceAccount)

		if err != nil {
			return err
		}
	}

	f.operation.SourceAccountID = sourceAccountID
	return nil
}

func (f *operationFactory) assignType() {
	f.operation.Type = strings.Replace(f.source.Body.Type.String(), "OperationType", "", 1)
}

func (f *operationFactory) assignSpecifics() error {
	var err error
	body := f.source.Body

	switch t := body.Type; t {
	case xdr.OperationTypeCreateAccount:
		f.assignCreateAccount(body.MustCreateAccountOp())
	case xdr.OperationTypePayment:
		err = f.assignPayment(body.MustPaymentOp())
	case xdr.OperationTypePathPaymentStrictReceive:
		err = f.assignPathPaymentStrictReceive(body.MustPathPaymentStrictReceiveOp())
	case xdr.OperationTypePathPaymentStrictSend:
		err = f.assignPathPaymentStrictSend(body.MustPathPaymentStrictSendOp())
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
		err = f.assignAccountMerge(body.MustDestination())
	case xdr.OperationTypeManageData:
		f.assignManageData(body.MustManageDataOp())
	case xdr.OperationTypeBumpSequence:
		f.assignBumpSequence(body.MustBumpSequenceOp())
	}

	return err
}

func (f *operationFactory) assignCreateAccount(o xdr.CreateAccountOp) {
	f.operation.SourceAmount = amount.String(o.StartingBalance)
	f.operation.DestinationAccountID = o.Destination.Address()
}

func (f *operationFactory) assignPayment(o xdr.PaymentOp) error {
	f.operation.SourceAmount = amount.String(o.Amount)

	var err error
	f.operation.DestinationAccountID, err = util.EncodeMuxedAccount(o.Destination)

	if err != nil {
		return err
	}

	f.operation.SourceAsset = NewAsset(&o.Asset)
	return nil
}

func (f *operationFactory) assignPathPaymentStrictReceive(o xdr.PathPaymentStrictReceiveOp) error {
	var err error
	f.operation.DestinationAccountID, err = util.EncodeMuxedAccount(o.Destination)

	if err != nil {
		return err
	}

	f.operation.DestinationAmount = amount.String(o.DestAmount)
	f.operation.DestinationAsset = NewAsset(&o.DestAsset)

	f.operation.SourceAmount = amount.String(o.SendMax)
	f.operation.SourceAsset = NewAsset(&o.SendAsset)

	f.operation.Path = make([]*Asset, len(o.Path))

	for i, a := range o.Path {
		f.operation.Path[i] = NewAsset(&a)
	}

	return nil
}

func (f *operationFactory) assignPathPaymentStrictSend(o xdr.PathPaymentStrictSendOp) error {
	var err error
	f.operation.DestinationAccountID, err = util.EncodeMuxedAccount(o.Destination)

	if err != nil {
		return err
	}

	f.operation.DestinationAmount = amount.String(o.DestMin)
	f.operation.DestinationAsset = NewAsset(&o.DestAsset)

	f.operation.SourceAmount = amount.String(o.SendAmount)
	f.operation.SourceAsset = NewAsset(&o.SendAsset)

	f.operation.Path = make([]*Asset, len(o.Path))
	for i, a := range o.Path {
		f.operation.Path[i] = NewAsset(&a)
	}

	return nil
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

	flags := xdr.TrustLineFlags(o.Authorize)

	if flags.IsAuthorized() {
		f.operation.Authorize = Full
	} else if flags.IsAuthorizedToMaintainLiabilitiesFlag() {
		f.operation.Authorize = MaintainLiabilities
	} else {
		f.operation.Authorize = None
	}
}

func (f *operationFactory) assignAccountMerge(d xdr.MuxedAccount) error {
	var err error
	f.operation.DestinationAccountID, err = util.EncodeMuxedAccount(d)

	if err != nil {
		return err
	}

	return nil
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
