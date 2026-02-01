package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres() *sql.DB {
	dsn := "host=localhost port=5432 user=postgres password=26120331de dbname=Booking sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
