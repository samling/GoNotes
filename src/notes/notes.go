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

	stmt.Close()

	return row
}

//func SetNoteTitle(row *sql.Row, db *sql.DB) error {
//	return nil
//}
//
//func SetNoteBody(row *sql.Row, db *sql.DB) error {
//	stmt, err := db.Prepare("insert into body values ? where id = ?")
//	check(err)
//	defer stmt.Close()
//
//	_, err = stmt.Exec(row.body, row.id)
//	check(err)
//
//	db.Commit()
//
//	stmt.Close()
//
//	return nil
//}

func DisplayNoteTitles(gui *gocui.Gui, rows *sql.Rows) {
	gui.Execute(func(gui *gocui.Gui) error {
		defer rows.Close()

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
