package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alifpay/croach/models"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
)

// use a single instance of Validate, it caches struct info
var valid = validator.New()

//parse body of http Request
func parseBody(r *http.Request, req interface{}) bool {
	if r.Body == nil {
		return false
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error", "app/common.go, parseBody, ioutil.ReadAll", err)
		return false
	}
	err = jsoniter.Unmarshal(b, req)
	if err != nil {
		log.Println("error", "app/common.go, parseBody, jsoniter.Unmarshal", err)
		return false
	}
	return true
}

//http Response writes in json format
func reply(w http.ResponseWriter, r models.Response) {
	b, err := jsoniter.Marshal(r)
	if err != nil {
		log.Println("error", "app/common.go, reply, jsoniter.Marshal", err)
		http.Error(w, "внутренняя ошибка, попробуйте позже", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(b)
}
