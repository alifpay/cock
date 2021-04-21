package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/alifpay/cock/models"
)

func GetTxns(ctx context.Context, f models.FilterTxn) (resp models.FilterResponse, err error) {

	sql := `SELECT id, service_name, external_ref, txn_type, account, currency, amount, fee, description,
	balance, regdate, stsdate, status, err_code, err FROM txns`

	where := make(map[string]clause)

	if f.ServiceName = strings.TrimSpace(f.ServiceName); len(f.ServiceName) > 0 {
		where["service_name"] = clause{Cond: "=", Param1: f.ServiceName}
	}

	if f.ExternalRef = strings.TrimSpace(f.ExternalRef); len(f.ExternalRef) > 0 {
		where["external_ref"] = clause{Cond: "=", Param1: f.ExternalRef}
	}

	if f.Account > 0 {
		where["account"] = clause{Cond: "=", Param1: f.Account}
	}

	where["regdate"] = clause{Cond: "BETWEEN", Param1: f.DateFrom, Param2: f.DateTo}

	whereSQL, args := sqlWhere(where, f.RowsOffset, f.RowsLimit, " ORDER BY id DESC", true)

	rows, err := conn.Query(ctx, sql+whereSQL, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	cs := make([]models.Txn, 0)
	for rows.Next() {
		c := models.Txn{}
		err = rows.Scan(
			&c.ID, &c.ServiceName, &c.ExternalRef, &c.TxnType, &c.Account, &c.Currency, &c.Amount, &c.Fee, &c.Description,
			&c.Balance, &c.Regdate, &c.Stsdate, &c.Status, &c.ErrCode, &c.Err)
		if err != nil {
			return
		}
		cs = append(cs, c)
	}
	resp.Result = cs
	//any filter changes offset must be assigned to 0
	if f.RowsOffset == 0 {
		w, a := sqlWhere(where, f.RowsOffset, f.RowsLimit, "", false)
		err = conn.QueryRow(ctx, "SELECT COUNT(id) FROM txns"+w, a...).Scan(&resp.Count)
	}
	return
}

func sqlWhere(where map[string]clause, offset, limit int, orderby string, offsetLimit bool) (string, []interface{}) {

	sql := " WHERE 1 = 1"
	ix := 1
	args := make([]interface{}, 0)
	for k, v := range where {
		switch v.Cond {
		case ">":
			args = append(args, v.Param1)
			sql += fmt.Sprintf(" AND %s > $%d", k, ix)
		case "<":
			args = append(args, v.Param1)
			sql += fmt.Sprintf(" AND %s < $%d", k, ix)
		case "=":
			args = append(args, v.Param1)
			sql += fmt.Sprintf(" AND %s = $%d", k, ix)
		case "LIKE":
			args = append(args, v.Param1)
			sql += fmt.Sprintf(" AND %s LIKE $%d", k, ix)
		case "BETWEEN":
			args = append(args, v.Param1)
			args = append(args, v.Param2)
			sql += fmt.Sprintf(" AND (%s BETWEEN $%d AND $%d)", k, ix, ix+1)
			ix++
		default:
			continue
		}
		ix++
	}

	if len(orderby) > 0 {
		sql += orderby
	}

	if offsetLimit {
		sql += fmt.Sprintf(" OFFSET $%d LIMIT $%d", ix, ix+1)
		args = append(args, offset)
		args = append(args, limit)
	}
	return sql, args
}

//Clause -
type clause struct {
	Cond   string
	Param1 interface{}
	Param2 interface{}
}
