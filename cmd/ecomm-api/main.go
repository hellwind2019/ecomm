package main

import (
	"log"

	"github.com/hellwind2019/ecomm/db"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connection established successfully")

	//st := storer.NewStorer(db.GetDB())
}
