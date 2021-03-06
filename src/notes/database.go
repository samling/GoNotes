package main

import (
	"database/sql"
	_ "fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	_ "os"
)

func InitDB(filepath string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", filepath)
	check(err)
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateTables(db *sqlx.DB) {
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

func GetNoteTitles(db *sqlx.DB) []interface{} {
	rows, err := db.Queryx("select Id, Title from Notes")
	check(err)

	titles := []Note{}

	for rows.Next() {
		scan, err := rows.SliceScan()
		check(err)
	}

	return nil
}

func GetNoteBody(db *sqlx.DB, currNote int) *sql.Row {
	stmt, err := db.Prepare("select body from Notes where id = ?")
	check(err)
	defer stmt.Close()

	row := stmt.QueryRow(currNote)
	check(err)

	stmt.Close()

	return row
}
