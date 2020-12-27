package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hculpan/kabtrek/galaxy"
	"github.com/hculpan/kabtrek/game"
	"github.com/hculpan/kabtrek/quadrant"
)

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	if err := game.InitScreen(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	g := galaxy.NewGalaxy(25, 5)
	g.Player = quadrant.NewEnterprise(game.RandomInt(10), game.RandomInt(10))
	g.Player.QuadrantX = game.RandomInt(8)
	g.Player.QuadrantY = game.RandomInt(8)
	g.SetActiveQuadrant(g.Player.QuadrantX, g.Player.QuadrantY)
	g.ScanNeighborQuadrants()

	loop(g)

	game.CloseScreen()
	os.Exit(0)
}

func handlePanic() {
	if r := recover(); r != nil {
		game.CloseScreen()
		fmt.Fprintf(os.Stderr, "PANIC: %s\n", r)
		os.Exit(1)
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

func loop(g *galaxy.Galaxy) {
	defer handlePanic()

	// Draw initial een
	g.Draw()

	// Setup event polling thread
	ch := game.PollForEvents()

	currTime := time.Now()
	updateCheck := 0
	paused := false

	for {
		q := g.GetActiveQuadrant()

		if q.IsPlayerDead() {
			playerDeadDisplay(g, ch)
			break
		} else if g.NumberOfKlingons == 0 {
			playerWinsDisplay(g, ch)
			break
		}

		if paused {
			game.EmitStr(0, 40, "PAUSED")
		} else {
			if time.Since(currTime).Seconds() > 0.5 {
				currTime = time.Now()
				q.UpdateTorpedoes()
				g.Draw()
				updateCheck++
			}
			if updateCheck > 1 {
				g.Update()
				g.Draw()
				updateCheck = 0
			}
		}

		if len(ch) > 0 {
			event := <-ch
			switch ev := event.(type) {
			case *tcell.EventResize:
				g.Draw()
			case *tcell.EventKey:
				if ev.Rune() == 32 && !paused {
					paused = true
				} else if ev.Rune() == 32 {
					paused = false
				} else {
					if q.UIState == quadrant.Normal && g.GameState == game.Quadrant {
						if ev.Key() == tcell.KeyESC {
							g.SetGameState(game.Quitting)
						}

						num := int(ev.Rune())
						if num >= 49 && num <= 57 {
							q.MoveObject(q.Player, num-48)
							g.Update()
							g.Draw()
						} else {
							switch num {
							case 's', 'S':
								q.UpdateState(quadrant.Shields)
							case 'w', 'W':
								q.UpdateState(quadrant.Weapons)
							case 'n', 'N':
								if g.Player.Shields > 0 {
									q.AddMessage("** Cannot go to warp with shields raised! **")
								} else {
									q.UpdateState(quadrant.NavigationX)
								}
							case 'l', 'L':
								g.SetGameState(game.LongRangeSensors)
							case 'c', 'C':
								g.SetGameState(game.GalaxyMap)
							}
						}
					} else if g.GameState == game.Quitting && (ev.Rune() == 'y' || ev.Rune() == 'Y') {
						return
					} else if g.GameState == game.Quitting && (ev.Rune() == 'n' || ev.Rune() == 'N') {
						g.SetGameState(game.Quadrant)
					} else {
						if ev.Key() == tcell.KeyESC {
							g.SetGameState(game.Quadrant)
							q.UpdateState(quadrant.Normal)
						}

						if q.AwaitingInput {
							num := int(ev.Rune())
							if num >= 48 && num <= 57 {
								q.CurrentInput += string(ev.Rune())
								g.Draw()
							} else if num == 8 && len(q.CurrentInput) > 0 { // backspace
								q.CurrentInput = q.CurrentInput[:len(q.CurrentInput)-1]
								g.Draw()
							} else if num == 13 && len(q.CurrentInput) > 0 { // Enter
								q.AcceptInput()
							}
						} else {
							q.HandleKeyForState(*ev)
						}
					}
				}
			}
			currTime = time.Now()
		}
	}

}

func playerWinsDisplay(g game.Game, ch chan tcell.Event) {
	game.ClearScreen()

	w, _ := game.Size()

	game.DrawBox(15, 2, w-15, 16)

	currentLine := 4
	msg := fmt.Sprintf("On Stardate %.1f, the Enterprise successfully destroyed the", g.GetStardate())
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)
	currentLine++
	msg = "last Klingon ship and won the war."
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)

	timeTaken := g.GetStardate() - 3100.1
	currentLine += 2
	msg = fmt.Sprintf("It took %.1f Stardates to win the war.", timeTaken)
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)
	currentLine++
	msg = fmt.Sprintf("This is an average of %.1f Stardates per enemy.", timeTaken/float64(g.GetStartingKlingons()))
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)
	currentLine++
	msg = "Starfleet Command congratulates you on your victory, and you"
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)
	currentLine++
	msg = "are hereby promoted to Admiral."
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)

	currentLine += 2
	msg = "Press ESC to quit"
	game.EmitStr(w/2-len(msg)/2, currentLine+2, msg)

	game.ShowScreen()

	waitForEsc(ch)
}

func playerDeadDisplay(g game.Game, ch chan tcell.Event) {
	game.ClearScreen()

	w, _ := game.Size()

	game.DrawBox(15, 2, w-15, 16)

	currentLine := 4
	msg := fmt.Sprintf("The Enterprise was destroyed on Stardate %.1f with all hands lost.", g.GetStardate())
	game.EmitStr(w/2-len(msg)/2, currentLine, msg)
	currentLine += 2
	numberDestroyed := g.GetStartingKlingons() - g.GetRemainingKlingons()
	switch {
	case numberDestroyed == 0:
		msg = "You suffered an ignominious defeat, failing to destroy even one enemy ship."
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
		msg = "Your humiliation is lessened only by the fact that your ship was destroyed,"
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
		msg = "lost with all of its crew - including you!"
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
		msg = "Your tactics will be studied through the ages as an example of"
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
		msg = "how not to conduct a war!"
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
	case numberDestroyed > 0:
		msg = fmt.Sprintf("While you managed to defeat %d enemies, the remaining %d Klingons destroyed", numberDestroyed, g.GetRemainingKlingons())
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
		msg = "all the Starbases and won the war!"
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
		msg = "Your defeat will go down in the annals of history!"
		game.EmitStr(w/2-len(msg)/2, currentLine, msg)
		currentLine++
	}
	msg = "Press ESC to quit"
	game.EmitStr(w/2-len(msg)/2, currentLine+2, msg)

	game.ShowScreen()

	waitForEsc(ch)
}
