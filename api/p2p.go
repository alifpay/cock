package api

import (
	"log"
	"net/http"

	"github.com/alifpay/cock/db"
	"github.com/alifpay/cock/models"
)

func p2p(w http.ResponseWriter, r *http.Request) {
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

	if req.Account2 == 0 || len(req.Currency2) != 3 {
		resp.Code = 400
		resp.Message = "invalid account2"
		reply(w, resp)
		return
	}

	if !req.Amount2.IsPositive() {
		resp.Code = 402
		resp.Message = "invalid amount"
		reply(w, resp)
		return
	}

	id, code, err := db.P2P(r.Context(), req)
	if err != nil {
		log.Println("db.P2P", err)
	}
	resp.Code = code
	resp.ID = id
	reply(w, resp)
}
