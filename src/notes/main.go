package main

import (
	_ "bufio"
	_ "database/sql"
	_ "fmt"
	"github.com/jroimartin/gocui"
	_ "github.com/mattn/go-sqlite3"
	_ "io"
	_ "io/ioutil"
	"log"
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

	// Set up our tables if they don't exist
	CreateTables(db)

	// Set up the GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	check(err)
	defer g.Close()

	// Set our layout for the GUI
	g.SetManagerFunc(layout)

	// Ctrl-C keybinding
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	// Check for loops in the main window
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	//n := &Note{"test title", "test value", "test tag"}
	//testTag := &Tag{"test name", "test member"}
	//n.Save(db)

	//n.Load(db)
	//fmt.Printf("%t\n", res)

	//res = testTag.Add(db)
	//fmt.Printf("%t\n", res)

	// List the contents of our Notes category (bucket)
	//categoryList(db)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if _, err := g.SetView("sidebar", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Open a database connection
		db := InitDB(dbpath)

		// List out all of our notes
		ListNotes(g, db)
	}

	if _, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
