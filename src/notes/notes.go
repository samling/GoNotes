package main

import (
	"database/sql"
	_ "encoding/json"
	"fmt"
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

func (n Note) Load(db *sql.DB) bool {
	rows, err := db.Query("select id, title, body from Notes")
	check(err)

	defer rows.Close()

	for rows.Next() {
		var id int
		var title string
		var body string

		err = rows.Scan(&id, &title, &body)
		check(err)
		fmt.Println(id, title, body)
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
