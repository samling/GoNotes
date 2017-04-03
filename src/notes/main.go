package main

import (
	_ "bufio"
	_ "database/sql"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "io"
	_ "io/ioutil"
	_ "log"
	_ "os"
)

const dbpath = "./database.db"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Open a database connection
	db := InitDB(dbpath)
	defer db.Close()
	CreateTables(db)

	// Close database connection after main() returns
	defer db.Close()

	n := &Note{"test title", "test value", "test tag"}
	//testTag := &Tag{"test name", "test member"}
	n.Save(db)

	n.Load(db)
	//fmt.Printf("%t\n", res)

	//res = testTag.Add(db)
	//fmt.Printf("%t\n", res)

	// List the contents of our Notes category (bucket)
	//categoryList(db)
}
