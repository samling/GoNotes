package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"github.com/jroimartin/gocui"
	_ "os"
)

type Note struct {
	title, body, tags string
}

type Tag struct {
	name, members string
}

func (n Note) Save(db *sql.DB) bool {
	tx, err := db.Begin()
	check(err)

	qry, err := tx.Prepare("insert into Notes(title, body) values (?, ?)")
	check(err)

	defer qry.Close()

	_, err = qry.Exec(n.title, n.body)
	check(err)

	tx.Commit()

	qry.Close()

	return true
}

func ListNotesStdOut(db *sql.DB) error {
	rows, err := db.Query("select id, title, body from Notes")
	check(err)

	for rows.Next() {
		var id int
		var title string
		var body string
		err = rows.Scan(&id, &title, &body)
		check(err)
		fmt.Println(id, title, body)
	}

	return nil
}

func GetNoteTitles(db *sql.DB) *sql.Rows {

	rows, err := db.Query("select id, title from Notes")
	check(err)

	return rows
}

func GetNoteBody(db *sql.DB, currNote int) *sql.Row {

	stmt, err := db.Prepare("select body from Notes where id = ?")
	check(err)
	defer stmt.Close()

	row := stmt.QueryRow(currNote)
	check(err)

	return row
}

func DisplayNoteTitles(gui *gocui.Gui, rows *sql.Rows) {
	gui.Execute(func(gui *gocui.Gui) error {
		s, err := gui.View("sidebar")
		check(err)
		s.Clear()

		for rows.Next() {
			var id int
			var title string

			err = rows.Scan(&id, &title)
			check(err)

			fmt.Fprintln(s, id, title)
		}

		err = rows.Err()
		check(err)

		rows.Close()
		return nil
	})
}

func DisplayNoteBody(gui *gocui.Gui, row *sql.Row) {
	gui.Execute(func(gui *gocui.Gui) error {
		m, err := gui.View("main")
		check(err)
		m.Clear()

		var body string

		err = row.Scan(&body)
		check(err)

		fmt.Fprintln(m, body)

		return nil
	})
}
