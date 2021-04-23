package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alifpay/cock/models"
)

var httpServer *http.Server //http server

func Run() {

	httpServer = &http.Server{
		Addr:              ":8095",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      40 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
		Handler:           routers(),
	}

	log.Println("web api server is running")
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}

func routers() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/acc", addAcc)
	mux.HandleFunc("/accs", getAccs)
	mux.HandleFunc("/credit", credit)
	mux.HandleFunc("/debit", debit)
	mux.HandleFunc("/p2p", p2p)
	mux.HandleFunc("/void", void)
	mux.HandleFunc("/balance", balance)
	mux.HandleFunc("/txns", getTxns)

	return mux
}

//Shutdown - graceful shutdown of http api server
func Shutdown(ctx context.Context) {
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

//Ping handler
func ping(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	reply(w, models.Response{Code: 200, Message: hostname})
}
