package db

import (
	"context"

	"github.com/alifpay/croach/models"
)

func GetBal(ctx context.Context, in models.Request) (acc models.Balance, err error) {
	err = conn.QueryRow(ctx, "SELECT external_ref, balance, name, regdate, status FROM accounts WHERE id = $1 AND currency = $2",
		in.Account, in.Currency).Scan(&acc.Ref, &acc.Balance, &acc.Name, &acc.Regdate, &acc.Status)
	return
}
