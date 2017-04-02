package main

import "github.com/boltdb/bolt"

func dbOpen() *bolt.DB {
	// Open a connection to the database and return a db object
	db, err := bolt.Open("notes.db", 0600, nil)
	check(err)
	return db
}
