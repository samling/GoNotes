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

	return true
}

func ListNotes(gui *gocui.Gui, db *sql.DB) bool {
	gui.Execute(func(gui *gocui.Gui) error {
		s, err := gui.View("sidebar")
		check(err)

		m, err := gui.View("main")
		check(err)

		s.Clear()
		m.Clear()

		rows, err := db.Query("select id, title, body from Notes")
		check(err)

		for rows.Next() {
			var id int
			var title string
			var body string

			err = rows.Scan(&id, &title, &body)
			check(err)

			fmt.Fprintf(s, "%s\n", title)
			fmt.Fprintln(m, body)
			return nil
		}

		err = rows.Err()
		check(err)

		rows.Close()
		return nil
	})

	return true
}
