package db

// XDR xdr.LedgerHeader `json:"-"`

// LedgerVersion int       `json:"version"`
// TotalCoins    int       `json:"total_coins"`
// FeePool       int       `json:"fee_pool"`
// InflationSeq  int       `json:"inflation_seq"`
// IDPool        int       `json:"id_pool"`
// BaseFee       int       `json:"base_fee"`
// BaseReserve   int       `json:"base_reserve"`
// MaxTxSetSize  int       `json:"max_tx_size"`
// CloseTime     time.Time `json:"close_time" db:"-"`

// func (l *Ledger) Unmarshal() {
// 	xdr.SafeUnmarshalBase64(l.Data, &l.XDR)

// 	l.LedgerVersion = int(l.XDR.LedgerVersion)
// 	l.TotalCoins = int(l.XDR.TotalCoins)
// 	l.FeePool = int(l.XDR.FeePool)
// 	l.InflationSeq = int(l.XDR.InflationSeq)
// 	l.IDPool = int(l.XDR.IdPool)
// 	l.BaseFee = int(l.XDR.BaseFee)
// 	l.BaseReserve = int(l.XDR.BaseReserve)
// 	l.MaxTxSetSize = int(l.XDR.MaxTxSetSize)
// 	l.CloseTime = time.Unix(l.CloseTimestamp, 0)
// }

// type Transaction struct {
// 	ID        string `db:"txid" json:"-"`
// 	LedgerSeq int    `db:"ledgerseq" json:"seq"`
// 	Index     int    `db:"txindex" json:"idx"`

// 	BodyXDR xdr.TransactionEnvelope `json:"-"`
// 	// ResultXDR xdr.TransactionResult   `json:"-"`
// 	// MetaXDR   xdr.TransactionMeta     `json:"-"`

// 	Body   string `db:"txbody" json:"-"`
// 	Result string `db:"txresult" json:"-"`
// 	Meta   string `db:"txmeta" json:"-"`
// }

// func (l *Transaction) Unmarshal() {
// 	xdr.SafeUnmarshalBase64(l.Body, &l.BodyXDR)
//}

// func GetTransactions(seq int) []Transaction {
// 	txs := []Transaction{}

// 	err := config.DB.Select(&txs, "SELECT * FROM txhistory WHERE ledgerseq = $1 ORDER BY txindex", seq)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return txs
// }
