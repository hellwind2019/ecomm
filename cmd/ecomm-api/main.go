package main

import (
	"log"

	"github.com/hellwind2019/ecomm/cmd/ecomm-api/handler"
	"github.com/hellwind2019/ecomm/cmd/ecomm-api/server"
	"github.com/hellwind2019/ecomm/cmd/ecomm-api/storer"
	"github.com/hellwind2019/ecomm/db"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully")

	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")

}
