package main

import (
	_ "bufio"
	_ "database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jroimartin/gocui"
	_ "github.com/mattn/go-sqlite3"
	_ "io"
	_ "io/ioutil"
	"log"
	_ "os"
)

const dbpath = "./database.db"

var (
	db       *sqlx.DB
	currNote int
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if s, err := g.SetView("sidebar", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		s.Title = "Notes"
		s.Editable = false
		s.Wrap = false
	}

	if m, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		m.Title = "Body"
		m.Editable = true
		m.Wrap = true
		m.Editor = &VimEditor{}

	}

	if _, err := g.SetCurrentView("sidebar"); err != nil {
		return err
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
	if err := gui.SetKeybinding("sidebar", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	if err := gui.SetKeybinding("msg", gocui.KeySpace, gocui.ModNone, delMsg); err != nil {
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

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
		if _, err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("sidebar"); err != nil {
		return err
	}
	return nil
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

	// Retrieve the list of notes
	//titles := GetNoteTitles(db)
	GetNoteTitles(db)
	body := GetNoteBody(db, currNote)

	// Set up the GUI
	gui, err := gocui.NewGui(gocui.OutputNormal)
	check(err)
	defer gui.Close()

	// Set our layout for the GUI
	gui.SetManagerFunc(layout)
	gui.InputEsc = true

	// Enable the cursor
	gui.Cursor = true

	// List out all of our notes
	//DisplayNoteTitles(gui, titles)

	// Display the body of our current note
	DisplayNoteBody(gui, body)

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
