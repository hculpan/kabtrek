package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
	"github.com/hculpan/kabtrek/galaxy"
	"github.com/hculpan/kabtrek/quadrant"
	"github.com/hculpan/kabtrek/util"
)

// Minimum width/height of console
const (
	MinWidth  = 80
	MinHeight = 25
)

func displayHelloWorld(s tcell.Screen) {
	w, h := s.Size()
	s.Clear()
	util.EmitStr(s, w/2-7, h/2, tcell.StyleDefault, "Hello, World!")
	util.EmitStr(s, w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")
	s.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	w, h := s.Size()
	if w < MinWidth || h < MinHeight {
		s.Fini()
		fmt.Fprintf(os.Stderr, "Console is only %d by %d characters in size\n", w, h)
		fmt.Fprintf(os.Stderr, "Console should be at least %d by %d characters in size\n", MinWidth, MinHeight)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)

	g := galaxy.NewGalaxy(25, 5)
	g.Player = quadrant.NewEnterprise(util.RandomInt(10), util.RandomInt(10))
	g.Player.QuadrantX = util.RandomInt(8)
	g.Player.QuadrantY = util.RandomInt(6)

	loop(s, g)

	s.Fini()
	os.Exit(0)
}

func handlePanic(scr tcell.Screen) {
	if r := recover(); r != nil {
		scr.Fini()
		fmt.Fprintf(os.Stderr, "PANIC: %s\n", r)
		os.Exit(1)
	}
}

func pollEvent(scr tcell.Screen, ch chan tcell.Event) {
	for {
		ch <- scr.PollEvent()
	}
}

func waitForEsc(ch chan tcell.Event) {
	t := true
	for t {
		if len(ch) > 0 {
			event := <-ch
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyESC {
					t = false
				}
			}
		}
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
}

func playerWinsDisplay(scr tcell.Screen, q *quadrant.Quadrant, ch chan tcell.Event) {
	scr.Sync()
	scr.Clear()

	w, _ := scr.Size()

	drawBox(scr, 15, 2, w-15, 16, tcell.StyleDefault)

	currentLine := 4
	msg := fmt.Sprintf("On Stardate %.1f, the Enterprise successfully destroyed the", q.Stardate)
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
	currentLine++
	msg = "last Klingon ship and won the war."
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)

	timeTaken := q.Stardate - 3100.1
	currentLine += 2
	msg = fmt.Sprintf("It took %.1f Stardates to win the war.", timeTaken)
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
	currentLine++
	msg = fmt.Sprintf("This is an average of %.1f Stardates per enemy.", timeTaken/float32(q.StartingNumberOfKlingons))
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
	currentLine++
	msg = "Starfleet Command congratulates you on your victory, and you"
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
	currentLine++
	msg = "are hereby promoted to Admiral."
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)

	currentLine += 2
	msg = "Press ESC to quit"
	util.EmitStr(scr, w/2-len(msg)/2, currentLine+2, tcell.StyleDefault, msg)

	scr.Show()

	waitForEsc(ch)
}

func playerDeadDisplay(scr tcell.Screen, q *quadrant.Quadrant, ch chan tcell.Event) {
	scr.Sync()
	scr.Clear()

	w, _ := scr.Size()

	drawBox(scr, 15, 2, w-15, 16, tcell.StyleDefault)

	currentLine := 4
	msg := fmt.Sprintf("The Enterprise was destroyed on Stardate %.1f with all hands lost.", q.Stardate)
	util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
	currentLine += 2
	numberDestroyed := q.StartingNumberOfKlingons - q.NumberOfKlingons
	switch {
	case numberDestroyed == 0:
		msg = "You suffered an ignominious defeat, failing to destroy even one enemy ship."
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
		msg = "Your humiliation is lessened only by the fact that your ship was destroyed,"
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
		msg = "lost with all of its crew - including you!"
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
		msg = "Your tactics will be studied through the ages as an example of"
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
		msg = "how not to conduct a war!"
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
	case numberDestroyed > 0:
		msg = fmt.Sprintf("While you managed to defeat %d enemies, the remaining %d Klingons destroyed", numberDestroyed, q.NumberOfKlingons)
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
		msg = "all the Starbases and won the war!"
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
		msg = "Your defeat will go down in the annals of history!"
		util.EmitStr(scr, w/2-len(msg)/2, currentLine, tcell.StyleDefault, msg)
		currentLine++
	}
	msg = "Press ESC to quit"
	util.EmitStr(scr, w/2-len(msg)/2, currentLine+2, tcell.StyleDefault, msg)
	scr.Show()

	waitForEsc(ch)
}

func loop(scr tcell.Screen, g *galaxy.Galaxy) {
	defer handlePanic(scr)

	q := &g.Quadrants[g.Player.QuadrantX][g.Player.QuadrantY]
	q.Player = g.Player
	for {
		x := util.RandomInt(10)
		y := util.RandomInt(10)
		if q.Objects[x][y] == nil {
			q.Player.X = x
			q.Player.Y = y
			q.Objects[x][y] = q.Player
			break
		}
	}

	// Draw initial screen
	q.Draw(scr)

	// Setup event polling thread
	ch := make(chan tcell.Event, 1)
	go pollEvent(scr, ch)

	currTime := time.Now()
	updateCheck := 0
	paused := false

	for {
		if q.IsPlayerDead() {
			playerDeadDisplay(scr, q, ch)
			break
		} else if g.NumberOfKlingons == 0 {
			playerWinsDisplay(scr, q, ch)
			break
		}

		if paused {
			util.EmitStr(scr, 0, 40, tcell.StyleDefault, "PAUSED")
		} else {
			if time.Since(currTime).Seconds() > 0.5 {
				currTime = time.Now()
				q.UpdateTorpedoes()
				q.Draw(scr)
				updateCheck++
			}
			if updateCheck > 1 {
				q.Update()
				q.Draw(scr)
				updateCheck = 0
			}
		}

		if len(ch) > 0 {
			event := <-ch
			switch ev := event.(type) {
			case *tcell.EventResize:
				q.Draw(scr)
			case *tcell.EventKey:
				if ev.Rune() == 32 && !paused {
					paused = true
				} else if ev.Rune() == 32 {
					paused = false
				} else {
					if q.UIState == quadrant.Normal {
						if ev.Key() == tcell.KeyESC {
							return
						}

						num := int(ev.Rune())
						if num >= 49 && num <= 57 {
							q.MoveObject(q.Player, num-48)
							q.Update()
							q.Draw(scr)
						} else {
							switch num {
							case 's', 'S':
								q.UpdateState(scr, quadrant.Shields)
							case 'w', 'W':
								q.UpdateState(scr, quadrant.Weapons)
							}
						}
					} else {
						if ev.Key() == tcell.KeyESC {
							q.UpdateState(scr, quadrant.Normal)
						}

						if q.AwaitingInput {
							num := int(ev.Rune())
							if num >= 48 && num <= 57 {
								q.CurrentInput += string(ev.Rune())
								q.Draw(scr)
							} else if num == 8 && len(q.CurrentInput) > 0 { // backspace
								q.CurrentInput = q.CurrentInput[:len(q.CurrentInput)-1]
								q.Draw(scr)
							} else if num == 13 && len(q.CurrentInput) > 0 { // Enter
								q.AcceptInput(scr)
							}
						} else {
							q.HandleKeyForState(scr, *ev)
						}
					}
				}
			}
			currTime = time.Now()
		}
	}

}
