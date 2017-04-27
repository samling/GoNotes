package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jroimartin/gocui"
	_ "os"
	"time"
)

type AutoIncr struct {
	ID      uint64
	Created time.Time
}

type Notebook struct {
	Notes []Note
	Tags  []Tag
}

type Note struct {
	AutoIncr
	SqlId int    `db:"id"`
	Title string `db:"title"`
	Body  string `db:"body"`
	Tags  []Tag
}

type Tag struct {
	AutoIncr
	Name    string
	Members []Note
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

func DisplayNoteTitles(gui *gocui.Gui, rows *sqlx.Rows) {
	gui.Execute(func(gui *gocui.Gui) error {
		defer rows.Close()

		//_, maxY := gui.Size()

		s, err := gui.View("sidebar")
		check(err)
		s.Clear()

		//if n, err := gui.SetView("noteId", -1, maxY-2, 10, maxY); err != nil {
		//	return err
		//}

		for rows.Next() {
			var id int
			var title string

			err = rows.Scan(&id, &title)
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
