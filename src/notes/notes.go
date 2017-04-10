package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"github.com/jroimartin/gocui"
	_ "os"
)

type Notebook struct {
	name string
	desc string
	note []Note
}

type Note struct {
	title, body, tags string
}

type Tag struct {
	name, members string
}

//func SetNoteTitle(note Note, db *sqlx.DB) error {
//	return nil
//}
//
//func SetNoteBody(note Note, db *sqlx.DB) error {
//	stmt, err := db.Prepare("insert into body values ? where id = ?")
//	check(err)
//	defer stmt.Close()
//
//	_, err = stmt.Exec(note.body, currNote)
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
			var title string

			err = rows.Scan(&title)
			check(err)

			fmt.Fprintln(s, title)
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
