package db

import (
	"context"
	"errors"
	"strings"

	"github.com/alifpay/croach/models"
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

func GetAccs(ctx context.Context, f models.FilterTxn) (resp models.FilterResponse, err error) {

	sql := `SELECT id, currency, external_ref, balance, name, regdate, status FROM accounts`

	where := make(map[string]clause)

	if f.ExternalRef = strings.TrimSpace(f.ExternalRef); len(f.ExternalRef) > 0 {
		where["external_ref"] = clause{Cond: "=", Param1: f.ExternalRef}
	}

	if f.Account > 0 {
		where["id"] = clause{Cond: "=", Param1: f.Account}
	}

	where["regdate"] = clause{Cond: "BETWEEN", Param1: f.DateFrom, Param2: f.DateTo}

	whereSQL, args := sqlWhere(where, f.RowsOffset, f.RowsLimit, " ORDER BY regdate DESC", true)

	rows, err := conn.Query(ctx, sql+whereSQL, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	cs := make([]models.Account, 0)
	for rows.Next() {
		c := models.Account{}
		err = rows.Scan(
			&c.Acc, &c.Currency, &c.Ref, &c.Balance, &c.Name, &c.Regdate, &c.Status)
		if err != nil {
			return
		}
		cs = append(cs, c)
	}
	resp.Result = cs
	//any filter changes offset must be assigned to 0
	if f.RowsOffset == 0 {
		w, a := sqlWhere(where, f.RowsOffset, f.RowsLimit, "", false)
		err = conn.QueryRow(ctx, "SELECT COUNT(id) FROM accounts"+w, a...).Scan(&resp.Count)
	}
	return
}
