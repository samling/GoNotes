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

type VimEditor struct {
	Insert bool
}

func (ve *VimEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ve.Insert {
		ve.InsertMode(v, key, ch, mod)
	} else {
		ve.NormalMode(v, key, ch, mod)
	}
}

func (ve *VimEditor) InsertMode(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case key == gocui.KeyEsc:
		ve.Insert = false
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		v.EditNewLine()
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
	// TODO: handle other keybindings...
}

func (ve *VimEditor) NormalMode(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch == 'i':
		ve.Insert = true
	case ch == 'j':
		v.MoveCursor(0, 1, false)
	case ch == 'k':
		v.MoveCursor(0, -1, false)
	case ch == 'h':
		v.MoveCursor(-1, 0, false)
	case ch == 'l':
		v.MoveCursor(1, 0, false)
	}
	// TODO: handle other keybindings...
}

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
		s.Editable = true
		s.Wrap = true
		s.Editor = &VimEditor{}
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
	rows := GetNoteTitles(db)
	row := GetNoteBody(db, currNote)

	// Set up the GUI
	gui, err := gocui.NewGui(gocui.OutputNormal)
	check(err)
	defer gui.Close()

	// Set our layout for the GUI
	gui.SetManagerFunc(layout)
	gui.InputEsc = true

	// Enable the cursor
	gui.Cursor = true

	// Set the foreground color
	gui.FgColor = gocui.ColorGreen
	gui.BgColor = gocui.ColorBlack

	// List out all of our notes
	DisplayNoteTitles(gui, rows)

	// Display the body of our current note
	DisplayNoteBody(gui, row)

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
