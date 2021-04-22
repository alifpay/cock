package api

import (
	"net/http"

	"github.com/alifpay/cock/db"
	"github.com/alifpay/cock/models"
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
