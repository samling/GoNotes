package main

import (
	"database/sql"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "os"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	check(err)
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTables(db *sql.DB) {
	qry := `
	create table if not exists Notes (
		id INTEGER NOT NULL PRIMARY KEY,
		title TEXT,
		body TEXT
	);
	`

	_, err := db.Exec(qry)
	check(err)
}
