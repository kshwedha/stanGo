package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func SqlCursor() *sql.DB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
