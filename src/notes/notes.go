package main

import (
	"fmt"
	"github.com/boltdb/bolt"
)

type Note struct {
	title, body, category string
}

type Tag struct {
	name, members string
}

type Category interface {
	Add(db *bolt.DB) bool
	//Modify(db *bolt.DB) bool
	//Remove(db *bolt.DB) bool
}

func (n Note) Add(db *bolt.DB) bool {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Notes"))
		err := b.Put([]byte(n.title), []byte(n.body))
		return err
	})

	check(err)

	return true
}

func (t Tag) Add(db *bolt.DB) bool {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Notes"))
		err := b.Put([]byte(t.name), []byte(t.members))
		return err
	})

	check(err)

	return true
}

func categoryList(b *bolt.DB) {
	// Open the "Notes" bucket and print its key/value pairs
	b.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Notes"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})

		return nil
	})
}
