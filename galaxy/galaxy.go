package galaxy

import (
	"fmt"

	"github.com/hculpan/kabtrek/game"
	"github.com/hculpan/kabtrek/quadrant"
)

// Galaxy contains info for, well, the galaxy
type Galaxy struct {
	Stardate                  float64
	StartingNumberOfKlingons  int
	StartingNumberOfStarbases int
	NumberOfKlingons          int
	NumberOfStarbases         int
	Quadrants                 [8][8]quadrant.Quadrant
	Player                    *quadrant.Enterprise

	ActiveQuadrantX int
	ActiveQuadrantY int
	GameState       int
}

// NewGalaxy create a whole new galaxy
func NewGalaxy(numKlingons, numStarbases int) *Galaxy {
	result := &Galaxy{
		Stardate:                  3700.1,
		StartingNumberOfKlingons:  numKlingons,
		StartingNumberOfStarbases: numStarbases,
		NumberOfKlingons:          numKlingons,
		NumberOfStarbases:         numStarbases,
		Quadrants:                 [8][8]quadrant.Quadrant{},
		GameState:                 game.Quadrant,
	}

	quadsGened := [8][8]bool{}
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			quadsGened[x][y] = false
		}
	}

	remainingKlingons := numKlingons
	remainingStarbases := numStarbases
	for remainingKlingons > 0 {
		// Pick number of klingons to place
		numKlingonsInQuadrant := selectKlingonsForQuadrant()
		if numKlingonsInQuadrant > remainingKlingons {
			numKlingonsInQuadrant = remainingKlingons
		}
		remainingKlingons -= numKlingonsInQuadrant

		for {
			x := game.RandomInt(8)
			y := game.RandomInt(8)
			if !quadsGened[x][y] && remainingStarbases > 0 && game.CheckPercent(10) {
				result.Quadrants[x][y] = *quadrant.NewQuadrant(result, x, y, numKlingonsInQuadrant, game.RandomInt(7), 1)
				remainingStarbases--
				quadsGened[x][y] = true
				break
			} else if !quadsGened[x][y] {
				result.Quadrants[x][y] = *quadrant.NewQuadrant(result, x, y, numKlingonsInQuadrant, game.RandomInt(7), 0)
				quadsGened[x][y] = true
				break
			}
		}

	}

	for remainingStarbases > 0 {
		x := game.RandomInt(8)
		y := game.RandomInt(8)
		if !quadsGened[x][y] {
			result.Quadrants[x][y] = *quadrant.NewQuadrant(result, x, y, 0, game.RandomInt(7), 1)
			remainingStarbases--
			quadsGened[x][y] = true
		}
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if !quadsGened[x][y] {
				result.Quadrants[x][y] = *quadrant.NewQuadrant(result, x, y, 0, game.RandomInt(7), 0)
			}
		}
	}

	return result
}

func selectKlingonsForQuadrant() int {
	n := game.GetPercent()
	switch {
	case n < 55:
		return 0
	case n < 75:
		return 1
	case n < 85:
		return 2
	case n < 92:
		return 3
	case n < 97:
		return 4
	default:
		return 5
	}
}

func (g *Galaxy) getQuadrant(x, y int) *quadrant.Quadrant {
	if x >= 0 && x < 8 && y >= 0 && y < 8 {
		return &g.Quadrants[x][y]
	}
	return nil
}

// ScanNeighborQuadrants reveals all quadrants
// adjacent to active quadrant
func (g *Galaxy) ScanNeighborQuadrants() {
	qx := g.ActiveQuadrantX
	qy := g.ActiveQuadrantY
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			q := g.getQuadrant(x+qx, y+qy)
			if q != nil {
				q.Scanned = true
			}
		}
	}
}

// Update runs the update cycle for the galaxy and
// any quadrants
func (g *Galaxy) Update() {
	if g.GameState == game.Quadrant {
		g.Stardate += 0.1

		g.GetActiveQuadrant().Update()
	}
}

// GetActiveQuadrant returns the active quadrant
func (g *Galaxy) GetActiveQuadrant() *quadrant.Quadrant {
	return &g.Quadrants[g.ActiveQuadrantX][g.ActiveQuadrantY]
}

// SetActiveQuadrant sets the active quadrant
func (g *Galaxy) SetActiveQuadrant(qx, qy int) {
	// First clear player from existing quadrant
	q := g.GetActiveQuadrant()
	if q != nil {
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				switch q.Objects[x][y].(type) {
				case quadrant.Player:
					q.Objects[x][y] = nil
					break
				}
			}
		}
	}

	// Now set new quadrant and add player to map
	g.ActiveQuadrantX, g.ActiveQuadrantY = qx, qy
	q = g.GetActiveQuadrant()
	if q == nil {
		panic(fmt.Sprintf("Unable to find quadrant %d, %d", qx, qy))
	}
	q.Scanned = true
	q.Player = g.Player
	for {
		x := game.RandomInt(10)
		y := game.RandomInt(10)
		if q.Objects[x][y] == nil {
			q.Player.X = x
			q.Player.Y = y
			q.Objects[x][y] = q.Player
			break
		}
	}
}

func (g *Galaxy) drawGalaxyMap() {
	x := 2
	y := 1
	game.EmitStr(x, y-1, "                     GALAXY MAP")
	game.EmitStr(x, y, " -------------------------------------------------")
	for yq := 0; yq < 8; yq++ {
		game.EmitStr(x, y+yq+1, "|                                                 |")
		for xq := 0; xq < 8; xq++ {
			lx := x + (xq * 6) + 1
			ly := y + yq + 1
			s := g.GetQuadrantSummary(xq, yq)
			if s.IsActive {
				game.EmitStr(lx, ly, fmt.Sprintf(" *%d%d%d*", s.Klingons, s.Starbases, s.Stars))
			} else if s.Scanned {
				game.EmitStr(lx, ly, fmt.Sprintf("  %d%d%d ", s.Klingons, s.Starbases, s.Stars))
			} else {
				game.EmitStr(lx, ly, "  ??? ")
			}
		}
	}
	game.EmitStr(x, y+9, " -------------------------------------------------")
	msg := fmt.Sprintf("STARDATE: %.1f     KLINGONS: %d     STARBASES: %d", g.Stardate, g.NumberOfKlingons, g.GetActiveQuadrant().NumberOfKlingons)
	game.EmitStr(2, y+11, msg)
}

func (g *Galaxy) drawLongRangeSensors() {
	g.ScanNeighborQuadrants()

	xloc := 25
	yloc := 8
	xq, yq := g.ActiveQuadrantX, g.ActiveQuadrantY
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			q := g.getQuadrant(xq+x, yq+y)
			game.EmitStr(xloc+(x*6), yloc+(y*4), "------")
			game.EmitStr(xloc+(x*6), yloc+(y*4)+1, "|     |")
			if q != nil {
				q.Scanned = true
				game.EmitStr(xloc+(x*6), yloc+(y*4)+2, fmt.Sprintf("| %d%d%d |", q.NumberOfKlingons, q.NumberOfStarbases, q.NumberOfStars))
			} else {
				game.EmitStr(xloc+(x*6), yloc+(y*4)+2, "| *** |")
			}
			game.EmitStr(xloc+(x*6), yloc+(y*4)+3, "|     |")
			game.EmitStr(xloc+(x*6), yloc+(y*4)+4, "------")
		}
	}
}

/*************************************
* methods to impliment Game interface
**************************************/

// NavigateTo the specified quadrant
func (g *Galaxy) NavigateTo(x, y int) {
	if x >= 0 && x <= 8 && y >= 0 && y <= 8 {
		g.SetActiveQuadrant(x, y)
		g.Update()
		g.Draw()
	}
}

func (g *Galaxy) quitting() {
	w, h := game.Size()
	msg := "Do you wish to quit (Y/N)?"
	game.EmitStr(w/2-len(msg)/2, h/2, msg)
}

// Draw draw's the quadrant
func (g *Galaxy) Draw() {
	game.ClearScreen()

	switch g.GameState {
	case game.GalaxyMap:
		g.drawGalaxyMap()
	case game.LongRangeSensors:
		g.drawLongRangeSensors()
	case game.Quitting:
		g.quitting()
	default:
		q := g.GetActiveQuadrant()
		q.DisplayQuadrant()
		q.DisplayStatus()
		q.DisplayState()
		q.DisplayMessages()
	}
	game.ShowScreen()
}

// GetQuadrantSummary returns the summary for the specified quadrant
func (g *Galaxy) GetQuadrantSummary(x, y int) *game.QuadrantSummary {
	if x >= 0 && x <= 8 && y >= 0 && y <= 8 {
		q := g.Quadrants[x][y]
		return &game.QuadrantSummary{
			X:         x,
			Y:         y,
			Klingons:  q.NumberOfKlingons,
			Starbases: q.NumberOfStarbases,
			Stars:     q.NumberOfStars,
			IsActive:  x == g.ActiveQuadrantX && y == g.ActiveQuadrantY,
			Scanned:   q.Scanned,
		}
	}

	return nil
}

// SetGameState sets the state of the game
func (g *Galaxy) SetGameState(state int) {
	g.GameState = state
	g.Draw()
}

// GetStartingKlingons returns the number of klingons
// game started with
func (g Galaxy) GetStartingKlingons() int {
	return g.StartingNumberOfKlingons
}

// GetRemainingKlingons gets the remaining klingons
func (g Galaxy) GetRemainingKlingons() int {
	return g.NumberOfKlingons
}

// KlingonDestroyed decrements the klingon counter
func (g *Galaxy) KlingonDestroyed() {
	g.NumberOfKlingons--
	if g.NumberOfKlingons < 0 {
		g.NumberOfKlingons = 0
	}
}

// GetStartingStarbases returns the number of klingons
// game started with
func (g *Galaxy) GetStartingStarbases() int {
	return g.StartingNumberOfStarbases
}

// GetRemainingStarbases returns the remaining starbases
func (g *Galaxy) GetRemainingStarbases() int {
	return g.NumberOfStarbases
}

// StarbaseDestroyed decrements the starbase counter
func (g *Galaxy) StarbaseDestroyed() {
	g.NumberOfStarbases--
	if g.NumberOfStarbases < 0 {
		g.NumberOfStarbases = 0
	}
}

// GetStardate gets the stardate
func (g *Galaxy) GetStardate() float64 {
	return g.Stardate
}
