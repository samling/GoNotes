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

func ListNotes(g *gocui.Gui, db *sql.DB) bool {
	g.Execute(func(g *gocui.Gui) error {
		s, err := g.View("sidebar")
		check(err)

		m, err := g.View("main")
		check(err)

		s.Clear()
		m.Clear()

		rows, err := db.Query("select id, title, body from Notes")
		check(err)
		defer rows.Close()

		for rows.Next() {
			var id int
			var title string
			var body string

			err = rows.Scan(&id, &title, &body)
			check(err)

			fmt.Fprintln(s, id, title)
			fmt.Fprintln(m, body)
			return nil
		}

		err = rows.Err()
		check(err)

		return nil
	})

	return true
}

func ViewList(db *sql.DB) bool {
	rows, err := db.Query("select id, title, body from Notes")
	check(err)

	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var body string

		err = rows.Scan(&id, &title, &body)
		fmt.Println(id, title, body, "\n")
		check(err)
	}
	err = rows.Err()
	check(err)

	return true
}

//
//func (n Note) Modify(noteName string, db *bolt.DB) bool {
//	err := db.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte("Notes"))
//		err := b.Put([]byte(n.title), []byte(n.body))
//		return err
//	})
//
//	check(err)
//
//	return true
//}
//
//func (t Tag) Add(db *bolt.DB) bool {
//	err := db.Update(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte("Notes"))
//		err := b.Put([]byte(t.name), []byte(t.members))
//		return err
//	})
//
//	check(err)
//
//	return true
//}
//
//func categoryList(b *bolt.DB) {
//	// Open the "Notes" bucket and print its key/value pairs
//	b.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte("Notes"))
//
//		b.ForEach(func(k, v []byte) error {
//			fmt.Printf("key=%s, value=%s\n", k, v)
//			return nil
//		})
//
//		return nil
//	})
//}
