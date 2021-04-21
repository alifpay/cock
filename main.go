package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alifpay/cock/api"
	"github.com/alifpay/cock/db"
)

func main() {
	connStr := "postgres://jack:secret@127.0.0.1:5432/bank?pool_max_conns=100"
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
