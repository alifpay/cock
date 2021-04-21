package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	Balance decimal.Decimal `json:"amount"`
	Name    string          `json:"name"`
	Status  string          `json:"status"`
	Ref     string          `json:"ref"`
	Regdate time.Time       `json:"regdate"`
}
