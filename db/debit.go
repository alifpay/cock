package db

import (
	"context"
	"errors"

	"github.com/alifpay/cock/models"
	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
)

func Debit(ctx context.Context, in models.Request) (id uint64, code int, err error) {

	str := `INSERT INTO txns(service_name, external_ref, txn_type, account, currency, amount, fee, description) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err = conn.QueryRow(ctx, str, in.ServiceName, in.ExternalRef, "debit", in.Account, in.Currency, in.Amount.Neg(), in.Fee, in.Description).Scan(&id)
	if err != nil {
		code = 504
		return
	}

	err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return debitTx(ctx, tx, in, id)
	})

	switch err {
	case nil:
		code = 200
		return
	case ErrBalance:
		code = 406
	case ErrNotActive:
		code = 405
	case pgx.ErrNoRows:
		code = 404
	default:
		code = 503 // db error
	}
	cmd, err := conn.Exec(ctx, "UPDATE txns SET stsdate = now(), status = 'failed', err_code = $1, err = $2 WHERE id = $3", code, err.Error(), id)
	if err != nil {
		return
	}
	if ra := cmd.RowsAffected(); ra == 0 {
		err = errors.New("no rows affected")
	}
	return
}

func debitTx(ctx context.Context, tx pgx.Tx, in models.Request, id uint64) error {
	// Read the balance.
	var (
		balance decimal.Decimal
		status  string
	)
	err := tx.QueryRow(ctx,
		"SELECT balance, status FROM accounts WHERE id = $1 AND currency = $2", in.Account, in.Currency).Scan(&balance, &status)
	if err != nil {
		return err
	}
	//check status
	if status != "active" {
		return ErrNotActive
	}

	if in.Amount.GreaterThan(balance) {
		return ErrBalance
	}

	// Perform the debit.
	cmd, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND currency = $3", in.Amount, in.Account, in.Currency)
	if err != nil {
		return err
	}
	if ra := cmd.RowsAffected(); ra == 0 {
		return errors.New("no rows affected")
	}
	//approve transaction
	cmd, err = tx.Exec(ctx, "UPDATE txns SET balance = $1, stsdate = now(), status = 'approved', err_code = 200 WHERE id = $2", balance, id)
	if err != nil {
		return err
	}
	if ra := cmd.RowsAffected(); ra == 0 {
		return errors.New("no rows affected")
	}
	return nil
}
