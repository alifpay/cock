package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// FilterResponse is a wrapper for any model with an additional field `TotalRecords`
type FilterResponse struct {
	Count  int         `json:"count"`
	Result interface{} `json:"result"`
}

type Txn struct {
	ID          uint64          `json:"id"`
	ServiceName string          `json:"service_name"`
	ExternalRef string          `json:"external_ref"`
	TxnType     string          `json:"txn_type"`
	Account     int             `json:"account"`
	Currency    string          `json:"currency"`
	Amount      decimal.Decimal `json:"amount"`
	Fee         decimal.Decimal `json:"fee"`
	Description string          `json:"description"`
	Balance     decimal.Decimal `json:"balance"`
	Regdate     time.Time       `json:"regdate"`
	Stsdate     time.Time       `json:"stsdate"`
	Status      string          `json:"status"`
	ErrCode     int             `json:"err_code"`
	Err         string          `json:"err"`
}

type FilterTxn struct {
	ServiceName string    `json:"service_name"`
	ExternalRef string    `json:"external_ref"`
	TxnType     string    `json:"txn_type"`
	Account     int       `json:"account"`
	Currency    string    `json:"currency"`
	DateFrom    time.Time `json:"date_from"`
	DateTo      time.Time `json:"date_to"`
	RowsLimit   int       `json:"rows_limit"`
	RowsOffset  int       `json:"rows_offset"`
}
