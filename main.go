package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"yatter-backend-go/app/config"
	"yatter-backend-go/app/dao"
	"yatter-backend-go/app/handler"
)

func main() {
	log.Fatalf("%+v", serve(context.Background()))
}

func serve(ctx context.Context) error {
	db, err := dao.NewDB(config.MySQLConfig())
	if err != nil {
		return err
	}
	defer db.Close()

	addr := ":" + strconv.Itoa(config.Port())
	log.Printf("Serve on http://%s", addr)

	return http.ListenAndServe(addr, handler.NewRouter(
		dao.NewAccount(db),
	))
}
