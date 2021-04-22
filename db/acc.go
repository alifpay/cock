package db

import (
	"context"
	"errors"

	"github.com/alifpay/cock/models"
)

func AddAcc(ctx context.Context, a models.Acc) error {
	str := "INSERT INTO accounts(id, currency, external_ref, name, status) VALUES($1, $2, $3, $4, 'active')"
	cmd, err := conn.Exec(ctx, str, a.Account, a.Currency, a.ExternalRef, a.Name)
	if err != nil {
		return err
	}

	if ra := cmd.RowsAffected(); ra == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
