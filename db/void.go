package db

import (
	"context"
	"errors"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
)

//void transaction
func Void(ctx context.Context, ServiceName, ExternalRef, description string) (code int, err error) {
	str := `INSERT INTO txns(service_name, external_ref, txn_type, account, currency, amount, fee, description) 
	SELECT service_name, external_ref, concat('void', txn_type), account, currency, amount*-1, fee, $1
	WHERE service_name = $2 AND external_ref = $3 AND txn_type IN ('debit', 'credit') AND status = 'approved'`

	cmd, err := conn.Exec(ctx, str, description, ServiceName, ExternalRef)
	if err != nil {
		code = 504
		return
	}
	if ra := cmd.RowsAffected(); ra == 0 {
		code = 504
		err = errors.New("no rows affected")
		return
	}

	str = `SELECT txn_type, account, currency, amount 
		   FROM txns service_name = $1 AND external_ref = $2 AND txn_type IN ('debit', 'credit') AND status = 'approved'`
	rows, err := conn.Query(ctx, str, ServiceName, ExternalRef)
	if err != nil {
		return
	}
	defer rows.Close()

	cs := make([]voidTxn, 0)
	for rows.Next() {
		c := voidTxn{}
		err = rows.Scan(&c.TxnType, &c.Account, &c.Currency, &c.Amount)
		if err != nil {
			return
		}
		cs = append(cs, c)
	}

	if ln := len(cs); ln == 0 {
		code = 504
		err = errors.New("no rows affected")
		return
	} else if ln == 2 {
		//void p2p
		err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return voidP2PTx(ctx, tx, cs, ServiceName, ExternalRef)
		})
	} else {
		//void credit or debit
		err = crdbpgx.ExecuteTx(ctx, conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
			return voidTx(ctx, tx, cs[0], ServiceName, ExternalRef)
		})
	}

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
	_, err = conn.Exec(ctx,
		`UPDATE txns SET stsdate = now(), status = 'failed', err_code = $1, err = $2 
		WHERE service_name = $3 AND external_ref = $4 AND txn_type IN ('voiddebit', 'voidcredit') AND status = 'pending'`,
		code, err.Error(), ServiceName, ExternalRef)
	return
}

func voidP2PTx(ctx context.Context, tx pgx.Tx, vs []voidTxn, ServiceName, ExternalRef string) error {
	var (
		balance decimal.Decimal
		status  string
	)
	// Read the balance.
	for _, v := range vs {
		if err := tx.QueryRow(ctx,
			"SELECT balance, status FROM accounts WHERE id = $1 AND currency = $2", v.Account, v.Currency).Scan(&balance, &status); err != nil {
			return err
		}
		//check status
		if status != "active" {
			return ErrNotActive
		}

		str := "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND currency = $3"
		if v.TxnType == "debit" {
			str = "UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND currency = $3"
		}
		if _, err := tx.Exec(ctx, str, v.Amount, v.Account, v.Currency); err != nil {
			return err
		}
		//approve transaction
		if _, err := tx.Exec(ctx,
			`UPDATE txns SET balance = $1, stsdate = now(), status = 'approved', err_code = 200 
			 WHERE service_name = $2 AND external_ref = $3 AND txn_type = $4`, balance, ServiceName, ExternalRef, v.TxnType); err != nil {
			return err
		}
	}
	return nil
}

func voidTx(ctx context.Context, tx pgx.Tx, v voidTxn, ServiceName, ExternalRef string) error {
	// Read the balance.
	var (
		balance decimal.Decimal
		status  string
	)
	if err := tx.QueryRow(ctx,
		"SELECT balance, status FROM accounts WHERE id = $1 AND currency = $2", v.Account, v.Currency).Scan(&balance, &status); err != nil {
		return err
	}
	//check status
	if status != "active" {
		return ErrNotActive
	}

	str := "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND currency = $3"
	if v.TxnType == "debit" {
		str = "UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND currency = $3"
	}
	if _, err := tx.Exec(ctx, str, v.Amount, v.Account, v.Currency); err != nil {
		return err
	}
	//approve transaction
	if _, err := tx.Exec(ctx,
		`UPDATE txns SET balance = $1, stsdate = now(), status = 'approved', err_code = 200 
		 WHERE service_name = $2 AND external_ref = $3 AND txn_type = $4`, balance, ServiceName, ExternalRef, v.TxnType); err != nil {
		return err
	}
	return nil
}

type voidTxn struct {
	TxnType  string
	Account  int
	Currency string
	Amount   decimal.Decimal
}
