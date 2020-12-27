package game

import (
	"fmt"

	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Minimum width/height of console
const (
	MinWidth  = 80
	MinHeight = 25
)

var scr tcell.Screen

// InitScreen initializes the screen
func InitScreen() error {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		return fmt.Errorf("%v", e)
	}
	scr = s

	if e := scr.Init(); e != nil {
		return e
	}

	w, h := scr.Size()
	if w < MinWidth || h < MinHeight {
		CloseScreen()
		return fmt.Errorf("console too small: found %d by %d, expected at least %d by %d", w, h, MinWidth, MinHeight)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	return nil
}

// CloseScreen releases the screen
func CloseScreen() {
	scr.Fini()
}

// ClearScreen clears the screen
func ClearScreen() {
	scr.Sync()
	scr.Clear()
}

// ShowScreen shows the screen
func ShowScreen() {
	scr.Show()
}

// Size returns the screen's size
func Size() (int, int) {
	return scr.Size()
}

// PollForEvents : Call this once, then check return channel
// for events periodically
func PollForEvents() chan tcell.Event {
	ch := make(chan tcell.Event, 1)
	go func() {
		for {
			ch <- scr.PollEvent()
		}
	}()
	return ch
}

// EmitStr will print a string to the screen
func EmitStr(x, y int, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		scr.SetContent(x, y, c, comb, tcell.StyleDefault)
		x += w
	}
}

// DrawBox draws a box
func DrawBox(x1, y1, x2, y2 int) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			scr.SetContent(col, row, ' ', nil, tcell.StyleDefault)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		scr.SetContent(col, y1, tcell.RuneHLine, nil, tcell.StyleDefault)
		scr.SetContent(col, y2, tcell.RuneHLine, nil, tcell.StyleDefault)
	}
	for row := y1 + 1; row < y2; row++ {
		scr.SetContent(x1, row, tcell.RuneVLine, nil, tcell.StyleDefault)
		scr.SetContent(x2, row, tcell.RuneVLine, nil, tcell.StyleDefault)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		scr.SetContent(x1, y1, tcell.RuneULCorner, nil, tcell.StyleDefault)
		scr.SetContent(x2, y1, tcell.RuneURCorner, nil, tcell.StyleDefault)
		scr.SetContent(x1, y2, tcell.RuneLLCorner, nil, tcell.StyleDefault)
		scr.SetContent(x2, y2, tcell.RuneLRCorner, nil, tcell.StyleDefault)
	}
}
