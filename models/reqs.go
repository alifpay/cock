package models

import (
	"github.com/shopspring/decimal"
)

type Request struct {
	ServiceName string          `json:"service_name" validate:"required,min=1,max=64"` //ServiceName partner service name - required
	ExternalRef string          `json:"external_ref" validate:"required,min=1,max=64"` //ExternalRef unique transaction id in partner system - required
	Account     int             `json:"account" validate:"required,len=8"`             //primary account  - required
	Currency    string          `json:"currency" validate:"required,len=3"`            //tjs, rub, usd, eur, uzs
	Amount      decimal.Decimal `json:"amount"`                                        //Amount 2.45 - required
	Account2    int             `json:"account2"`                                      //receiver account - optional, required in account to account operation
	Currency2   string          `json:"currency2"`                                     //tjs, rub, usd, eur, uzs
	Amount2     decimal.Decimal `json:"amount2"`                                       //Amount 2.45 - required
	Description string          `json:"description"`
	Fee         decimal.Decimal `json:"fee"`
}

type VoidReq struct {
	ServiceName string `json:"service_name" validate:"required,min=1,max=64"` //ServiceName partner service name - required
	ExternalRef string `json:"external_ref" validate:"required,min=1,max=64"` //ExternalRef unique transaction id in partner system - required
	Description string `json:"description" validate:"required,min=1,max=255"`
}
