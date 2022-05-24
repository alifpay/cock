package api

import (
	"net/http"

	"github.com/alifpay/croach/db"
	"github.com/alifpay/croach/models"
)

func addAcc(w http.ResponseWriter, r *http.Request) {
	var (
		req  models.Acc
		resp models.Response
	)
	if !parseBody(r, &req) {
		resp.Code = 400
		reply(w, resp)
		return
	}

	err := valid.StructCtx(r.Context(), req)
	if err != nil {
		resp.Code = 400
		resp.Message = err.Error()
		reply(w, resp)
		return
	}

	err = db.AddAcc(r.Context(), req)
	if err != nil {
		resp.Code = 503
		reply(w, resp)
		return
	}
	resp.Code = 200
	reply(w, resp)
}

func getAccs(w http.ResponseWriter, r *http.Request) {

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

	payload, err := db.GetAccs(r.Context(), f)
	if err != nil {
		resp.Code = 503
		reply(w, resp)
		return
	}

	resp.Code = 200
	resp.Payload = payload
	reply(w, resp)
}
