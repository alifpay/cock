package api

import (
	"net/http"

	"github.com/alifpay/croach/db"
	"github.com/alifpay/croach/models"
)

func balance(w http.ResponseWriter, r *http.Request) {
	var (
		req  models.Request
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

	payload, err := db.GetBal(r.Context(), req)
	if err != nil {
		if db.IsNotFound(err) {
			resp.Code = 404
		} else {
			resp.Code = 503
		}
		reply(w, resp)
		return
	}
	resp.Code = 200
	resp.Payload = payload
	reply(w, resp)
}
