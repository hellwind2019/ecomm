package main

import (
	"log"

	"github.com/hellwind2019/ecomm/cmd/ecomm-api/handler"
	"github.com/hellwind2019/ecomm/cmd/ecomm-api/server"
	"github.com/hellwind2019/ecomm/cmd/ecomm-api/storer"
	"github.com/hellwind2019/ecomm/db"
	"github.com/ianschenck/envflag"
)

const minSecretKeyLength = 32

func main() {
	var secretKey = envflag.String("SECRET_KEY", "0123456789012345678901234567890123456789019", "Secret key for JWT signing")
	if len(*secretKey) < minSecretKeyLength {
		log.Fatalf("Secret key must be at least %d characters long", minSecretKeyLength)
	}
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully")

	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv, *secretKey)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")

}
