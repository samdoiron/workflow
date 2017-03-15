package workflow

import (
	"log"

	"github.com/jmoiron/sqlx"

	// For the postgres driver
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func getConnection() *sqlx.DB {
	return db
}

func init() {
	newDB, err := sqlx.Open("postgres",
		"user=workflow dbname=workflow sslmode=disable")
	if err != nil {
		log.Fatal("[ERROR] failed to connect to database")
	}
	db = newDB
}
