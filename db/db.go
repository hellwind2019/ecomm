package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase() (*Database, error) {
	database, err := sqlx.Open("mysql", "root:admin@tcp(localhost:3306)/ecomm?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Database{db: database}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
