package main

import (
	_ "bufio"
	"fmt"
	"github.com/boltdb/bolt"
	_ "io"
	_ "io/ioutil"
	_ "log"
	_ "os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	// Open a database connection
	db := dbOpen()

	// Create base tables
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Notes"))
		check(err)
		return nil
	})

	// Close database connection after init() returns
	defer db.Close()

}

func main() {
	// Open a database connection
	db := dbOpen()

	testNote := &Note{"test title", "test value", "test category"}
	testTag := &Tag{"test name", "test member"}
	res := testNote.Add(db)
	fmt.Printf("%t\n", res)

	res = testTag.Add(db)
	fmt.Printf("%t\n", res)

	// List the contents of our Notes category (bucket)
	categoryList(db)

	// Close database connection after main() returns
	defer db.Close()
}
