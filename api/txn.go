package api

import (
	"net/http"

	"github.com/alifpay/cock/db"
	"github.com/alifpay/cock/models"
)

func getTxns(w http.ResponseWriter, r *http.Request) {

	var (
		f    models.FilterTxn
		resp models.Response
	)
	if !parseBody(r, &f) {
		resp.Code = 400
		reply(w, resp)
		return
	}

	if f.RowsLimit == 0 || f.RowsLimit > 1000 {
		resp.Code = 400
		reply(w, resp)
		return
	}

	if f.DateFrom.IsZero() || f.DateTo.IsZero() {
		resp.Code = 400
		reply(w, resp)
		return
	}

	payload, err := db.GetTxns(r.Context(), f)
	if err != nil {
		resp.Code = 503
		reply(w, resp)
		return
	}

	resp.Code = 200
	resp.Payload = payload
	reply(w, resp)
}
