package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alifpay/croach/api"
	"github.com/alifpay/croach/db"
	"github.com/shopspring/decimal"
)

func main() {
	decimal.MarshalJSONWithoutQuotes = true
	host := os.Getenv("DBHOST")
	if host == "" {
		host = "127.0.0.1"
	}
	connStr := "postgres://root@" + host + ":26257/bank?sslmode=disable&pool_max_conns=100"
	err := db.Connect(connStr)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancelFun := context.WithCancel(context.Background())
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-quitCh
		cancelFun()
		api.Shutdown(ctx)
	}()

	api.Run()
}
