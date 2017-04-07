package main

import (
	_ "bufio"
	"database/sql"
	_ "fmt"
	"github.com/jroimartin/gocui"
	_ "github.com/mattn/go-sqlite3"
	_ "io"
	_ "io/ioutil"
	"log"
	_ "os"
)

const dbpath = "./database.db"

var (
	db       *sql.DB
	currNote int
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Set the default current note
	// TODO: Make this save the last open note
	currNote = 1

	// Open a database connection
	db = InitDB(dbpath)
	defer db.Close()

	// Create our tables if they don't exist
	CreateTables(db)

	// Set up the GUI
	gui, err := gocui.NewGui(gocui.OutputNormal)
	check(err)
	defer gui.Close()

	// Set our layout for the GUI
	gui.SetManagerFunc(layout)

	// Enable the cursor
	gui.Cursor = true

	// List out all of our notes
	ListNotes(gui, db)

	// Ctrl-C keybinding
	if err := keybindings(gui); err != nil {
		log.Panicln(err)
	}

	// Check for loops in the main window
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	//n := &Note{"second title", "second value", "second tag"}
	//testTag := &Tag{"test name", "test member"}
	//n.Save(db)

	//n.Load(db)
	//fmt.Printf("%t\n", res)

	//res = testTag.Add(db)
	//fmt.Printf("%t\n", res)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if _, err := g.SetView("sidebar", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

	}

	if _, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func keybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("sidebar", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := gui.SetKeybinding("sidebar", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
