package api

import (
	"log"
	"net/http"

	"github.com/alifpay/cock/db"
	"github.com/alifpay/cock/models"
)

func void(w http.ResponseWriter, r *http.Request) {
	var (
		req  models.VoidReq
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

	code, err := db.Void(r.Context(), req.ServiceName, req.ExternalRef, req.Description)
	if err != nil {
		log.Println("db.Credit", err)
	}
	resp.Code = code
	reply(w, resp)
}
