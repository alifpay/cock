package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Balance decimal.Decimal `json:"balance"`
	Name    string          `json:"name"`
	Status  string          `json:"status"`
	Ref     string          `json:"ref"`
	Regdate time.Time       `json:"regdate"`
}

type Account struct {
	Acc      int             `json:"account"`
	Currency string          `json:"currency"`
	Balance  decimal.Decimal `json:"balance"`
	Name     string          `json:"name"`
	Status   string          `json:"status"`
	Ref      string          `json:"ref"`
	Regdate  time.Time       `json:"regdate"`
}
