package api

import (
	"log"
	"net/http"

	"github.com/alifpay/cock/db"
	"github.com/alifpay/cock/models"
)

func credit(w http.ResponseWriter, r *http.Request) {
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

	if !req.Amount.IsPositive() {
		resp.Code = 402
		resp.Message = "invalid amount"
		reply(w, resp)
		return
	}

	id, code, err := db.Credit(r.Context(), req)
	if err != nil {
		log.Println("db.Credit", err)
	}
	resp.Code = code
	resp.ID = id
	reply(w, resp)
}
