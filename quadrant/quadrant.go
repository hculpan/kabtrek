package quadrant

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/hculpan/kabtrek/game"
)

// Direction constants
const (
	_ = iota
	Dir1
	Dir2
	Dir3
	Dir4
	Dir5
	Dir6
	Dir7
	Dir8
	Dir9
)

// State of interface - which option currently selected
const (
	Normal = iota
	NavigationX
	NavigationY
	Weapons
	Shields
	Sensors
	Computer

	WeaponsPhasers
	WeaponsTorpedoes
)

// Game constants
const (
	EnergyToMove  = 10
	TorpedoDamage = 500
)

// Message type for ui messages
type Message struct {
	Text     string
	Stardate float64
}

// Quadrant : All the information related to a single Quadrant
type Quadrant struct {
	X                         int
	Y                         int
	Objects                   [10][10]Object
	Player                    *Enterprise
	StartingNumberOfKlingons  int
	NumberOfKlingons          int
	StartingNumberOfStarbases int
	NumberOfStarbases         int
	NumberOfStars             int
	Scanned                   bool
	Game                      game.Game

	UIState       int
	AwaitingInput bool
	CurrentInput  string

	Messages []Message

	// Private variables
	blinkRed     int
	torpedoes    [10][10]*Torpedo
	destinationX int
}

// NewQuadrant creates a new quadrant, populated with items
func NewQuadrant(parentGame game.Game, x int, y int, numKlingons int, numStars int, numBases int) *Quadrant {
	result := &Quadrant{
		Game:                      parentGame,
		X:                         x,
		Y:                         y,
		Objects:                   [10][10]Object{},
		Player:                    nil,
		NumberOfKlingons:          numKlingons,
		StartingNumberOfKlingons:  numKlingons,
		NumberOfStarbases:         numBases,
		StartingNumberOfStarbases: numBases,
		NumberOfStars:             numStars,
		UIState:                   0,
		AwaitingInput:             false,
		Scanned:                   false,
		CurrentInput:              "",
		blinkRed:                  0,
		torpedoes:                 [10][10]*Torpedo{},
	}

	result.blinkRed = 0
	result.UIState = Normal
	result.AwaitingInput = false
	result.CurrentInput = ""

	result.NumberOfKlingons = numKlingons
	klingonsToPlace := numKlingons
	for klingonsToPlace > 0 {
		xloc := game.RandomInt(10)
		yloc := game.RandomInt(10)
		if result.Objects[xloc][yloc] == nil {
			result.Objects[xloc][yloc] = NewKlingon(xloc, yloc)
			klingonsToPlace--
		}
	}

	starsToPlace := numStars
	for starsToPlace > 0 {
		xloc := game.RandomInt(10)
		yloc := game.RandomInt(10)
		if result.Objects[xloc][yloc] == nil {
			result.Objects[xloc][yloc] = &Star{X: xloc, Y: yloc}
			starsToPlace--
		}
	}

	basesToPlace := numBases
	for basesToPlace > 0 {
		xloc := game.RandomInt(10)
		yloc := game.RandomInt(10)
		if result.Objects[xloc][yloc] == nil {
			result.Objects[xloc][yloc] = NewStarbase(xloc, yloc)
			basesToPlace--
		}
	}

	return result
}

func (q *Quadrant) isBaseAt(x int, y int) bool {
	if x < 0 || x > 9 || y < 0 || y > 9 || q.Objects[x][y] == nil {
		return false
	} else {
		switch q.Objects[x][y].(type) {
		case *Starbase:
			return true
		default:
			return false
		}
	}
}

func (q *Quadrant) playerDockedAtBase() bool {
	x, y := q.Player.Location()
	if q.Player.Shields == 0 &&
		(q.isBaseAt(x, y-1) ||
			q.isBaseAt(x-1, y) ||
			q.isBaseAt(x+1, y) ||
			q.isBaseAt(x, y+1)) {
		return true
	}
	return false
}

// UpdateState changes the current state of the UI
func (q *Quadrant) UpdateState(newState int) {
	q.UIState = newState
	if newState != Normal && newState != Weapons && newState != NavigationX && newState != NavigationY {
		q.AwaitingInput = true
	} else {
		q.AwaitingInput = false
	}
	q.CurrentInput = ""
	q.Game.Draw()
}

// UpdateTorpedoes moves the torpedo along it's path
func (q *Quadrant) UpdateTorpedoes() {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			q.updateTorpedoAt(x, y)
		}
	}
}

// IsPlayerDead checks if game is over
func (q *Quadrant) IsPlayerDead() bool {
	return q.Player.Energy <= 0
}

func (q *Quadrant) damageObjectAt(x int, y int, damage int, deiptor string) {
	if q.Objects[x][y] == nil {
		return
	}

	q.Objects[x][y].TakeDamage(TorpedoDamage)
	q.AddMessage(fmt.Sprintf("%s at %d, %d took %d damage from a %s", q.Objects[x][y].Name(), x, y, damage, deiptor))

	if q.Objects[x][y].GetShields() <= 0 {
		q.AddMessage(fmt.Sprintf("%s at %d, %d destroyed!", q.Objects[x][y].Name(), x, y))
		switch q.Objects[x][y].(type) {
		case *Klingon:
			q.NumberOfKlingons--
			q.Game.KlingonDestroyed()
		}
		q.Objects[x][y] = nil
	}
}

func (q *Quadrant) updateTorpedoAt(ox int, oy int) {
	t := q.torpedoes[ox][oy]
	if t == nil {
		return
	}

	// Pick new location
	x, y := newLocation(ox, oy, t.Direction)

	// Check if goes off the board
	if x < 0 || x > 9 || y < 0 || y > 9 {
		q.torpedoes[ox][oy] = nil
	} else if q.Objects[x][y] != nil { // Has it hit anything?
		q.damageObjectAt(x, y, TorpedoDamage, "torpedo")
		q.torpedoes[ox][oy] = nil
	} else {
		q.torpedoes[x][y] = t
		q.torpedoes[ox][oy] = nil
	}
}

func (q *Quadrant) klingonFireTorpedo(k *Klingon, dir int) {
	ox, oy := k.Location()
	q.AddMessage(fmt.Sprintf("Klingon at %d, %d is firing a torpedo!", ox, oy))
	q.torpedoes[ox][oy] = &Torpedo{X: ox, Y: oy, Direction: dir}
	q.updateTorpedoAt(ox, oy)
}

func (q *Quadrant) klingonAction(k *Klingon) {
	x, y := k.Location()
	if game.CheckPercent(66) {
		if game.CheckPercent(25) {
			// Fire torpedo
			dx := x - q.Player.X
			dy := y - q.Player.Y
			if dx < 0 && dy < 0 {
				q.klingonFireTorpedo(k, 3)
			} else if dx == 0 && dy < 0 {
				q.klingonFireTorpedo(k, 2)
			} else if dx > 0 && dy < 0 {
				q.klingonFireTorpedo(k, 1)
			} else if dx < 0 && dy == 0 {
				q.klingonFireTorpedo(k, 6)
			} else if dx > 0 && dy == 0 {
				q.klingonFireTorpedo(k, 4)
			} else if dx > 0 && dy > 0 {
				q.klingonFireTorpedo(k, 7)
			} else if dx == 0 && dy > 0 {
				q.klingonFireTorpedo(k, 8)
			} else if dx < 0 && dy > 0 {
				q.klingonFireTorpedo(k, 9)
			}
		} else {
			// Move
			var o MoveableObject = q.Objects[x][y].(MoveableObject)
			q.MoveObject(o, game.RandomInt(9)+1)
		}
	}
}

// AddMessage adds a message to the messages display
// Will be removed in 5 turns (stardate + 0.5)
func (q *Quadrant) AddMessage(t string) {
	q.Messages = append(q.Messages, Message{Text: t, Stardate: q.Game.GetStardate()})
}

// Update processes the next turn for the quadrant
func (q *Quadrant) Update() {
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			switch q.Objects[x][y].(type) {
			case *Klingon:
				q.klingonAction(q.Objects[x][y].(*Klingon))
			}
		}
	}

	if q.playerDockedAtBase() {
		q.Player.Energy += game.EnterpriseMaxEnergy * 0.25
		if q.Player.Energy > game.EnterpriseMaxEnergy {
			q.Player.Energy = game.EnterpriseMaxEnergy
		}
		q.Player.Torpedoes += game.EnterpriseMaxTorpedoes * 0.25
		if q.Player.Torpedoes > game.EnterpriseMaxTorpedoes {
			q.Player.Torpedoes = game.EnterpriseMaxTorpedoes
		}
	}
}

// DisplayQuadrant draws the Quadrant map
func (q *Quadrant) DisplayQuadrant() {
	QuadrantStr := fmt.Sprintf("Quadrant : %d, %d", q.X+1, q.Y+1)
	game.EmitStr(22-(len(QuadrantStr)/2), 0, QuadrantStr)
	game.EmitStr(3, 1, "=---=---=---=---=---=---=---=---=---=---")
	for i := 0; i < 10; i++ {
		q.displayQuadrantLine(i)
	}
	game.EmitStr(3, 12, "=-1-=-2-=-3-=-4-=-5-=-6-=-7-=-8-=-9-=-10")
}

// DisplayMessages displays all the messages in order,
// removing any greater than 0.5 Stardates old
func (q *Quadrant) DisplayMessages() {
	for i, t := range q.Messages {
		game.EmitStr(1, 16+i, fmt.Sprintf("Stardate %.1f: %s", t.Stardate, t.Text))
	}

	// Remove older messages
	for len(q.Messages) > 0 {
		if q.Messages[0].Stardate < q.Game.GetStardate()-0.5 {
			q.Messages = q.Messages[1:]
		} else {
			break
		}
	}
}

func (q *Quadrant) displaySector(x int, y int) {
	if q.Objects[x][y] != nil {
		objStr := ""
		switch q.Objects[x][y].(type) {
		case Player:
			objStr = "-E-"
		case *Klingon:
			objStr = "-K-"
		case *Star:
			objStr = " * "
		case *Starbase:
			objStr = ">B<"
		}
		game.EmitStr(x*4+4, y+2, objStr)
	} else if q.torpedoes[x][y] != nil {
		game.EmitStr(x*4+4, y+2, " @ ")
	}

}

func (q *Quadrant) displayQuadrantLine(row int) {
	game.EmitStr(0, row+2, fmt.Sprintf("%2d|                                        |", row+1))

	for i := 0; i < 10; i++ {
		q.displaySector(i, row)
	}
}

// HandleKeyForState handles the current key for the various ui states
func (q *Quadrant) HandleKeyForState(key tcell.EventKey) {
	switch q.UIState {
	case Weapons:
		switch int(key.Rune()) {
		case 'p', 'P':
			q.UpdateState(WeaponsPhasers)
			q.Game.Draw()
		case 't', 'T':
			q.UpdateState(WeaponsTorpedoes)
			q.Game.Draw()
		}
	case NavigationX:
		num := int(key.Rune())
		if num >= 49 && num <= 56 {
			q.destinationX = num - 48
			q.UpdateState(NavigationY)
		}
	case NavigationY:
		num := int(key.Rune())
		if num >= 49 && num <= 56 {
			q.UpdateState(Normal)
			d := game.Distance(q.X, q.Y, q.destinationX-1, num-49)
			if q.Player.Energy < int(d*100) {
				q.AddMessage("You do not have enough energy for that trip")
			} else {
				q.Player.Energy -= int(d * 100)
				q.Game.NavigateTo(q.destinationX-1, num-49)
			}
		}
	}
}

// AcceptInput accepts whatever the player has typed for input
func (q *Quadrant) AcceptInput() {
	value, _ := strconv.Atoi(q.CurrentInput)
	switch q.UIState {
	case Shields:
		if value > q.Player.Energy+q.Player.Shields {
			q.Player.Shields = q.Player.Energy + q.Player.Shields
			q.Player.Energy = 0
		} else {
			q.Player.Energy = (q.Player.Energy + q.Player.Shields) - value
			q.Player.Shields = value
		}
	case WeaponsTorpedoes:
		if q.Player.Torpedoes > 1 && value >= 1 && value <= 9 && value != 5 {
			q.Player.Torpedoes--
			t := &Torpedo{X: q.Player.X, Y: q.Player.Y, Direction: value}
			q.AddMessage("Torpedo fired!")
			q.torpedoes[q.Player.X][q.Player.Y] = t
			q.updateTorpedoAt(q.Player.X, q.Player.Y)
		}
	}
	q.UpdateState(Normal)
	q.Game.Draw()
}

// DisplayState renders the bottom display
// based on the state of the UI
func (q *Quadrant) DisplayState() {
	switch q.UIState {
	case Normal:
		game.EmitStr(1, 14, "(N)avigation   (W)eapons   (S)hields  (L)ong-Range Sensors  Ship's (C)omputer")
	case Shields:
		game.EmitStr(1, 14, "Set energy for shields: ")
		q.displayInput(25, 14)
	case Weapons:
		game.EmitStr(1, 14, "(P)hasers or Photon (T)orpedoes")
	case WeaponsTorpedoes:
		game.EmitStr(1, 14, "Direction:")
		q.displayInput(12, 14)
	case NavigationX:
		game.EmitStr(1, 14, "Destination Quadrant X:")
	case NavigationY:
		game.EmitStr(1, 14, "Destination Quadrant Y:")
	}
}

func (q *Quadrant) displayInput(col int, row int) {
	game.EmitStr(col, row, q.CurrentInput)
	if q.blinkRed%2 == 0 {
		game.EmitStr(col+len(q.CurrentInput), row, "_")
	}
}

// DisplayStatus draws the status of the quadrant (the stuff to the right of the map)
func (q *Quadrant) DisplayStatus() {
	q.blinkRed++
	game.EmitStr(49, 3, fmt.Sprintf("STARDATE:         %.1f", q.Game.GetStardate()))
	game.EmitStr(49, 4, fmt.Sprintf("SECTOR:           %d,%d", q.Player.X, q.Player.Y))

	if q.playerDockedAtBase() {
		game.EmitStr(49, 5, "CONDITION:        DOCKED")
	} else if q.NumberOfKlingons > 0 && q.blinkRed%2 == 1 {
		game.EmitStr(49, 5, "CONDITION:        RED")
	} else if q.NumberOfKlingons > 0 {
		game.EmitStr(49, 5, "CONDITION: ")
	} else {
		game.EmitStr(49, 5, "CONDITION:        GREEN")
	}

	game.EmitStr(49, 6, fmt.Sprintf("SHIELDS:          %d", q.Player.Shields))
	game.EmitStr(49, 7, fmt.Sprintf("ENERGY:           %d", q.Player.Energy))
	game.EmitStr(49, 8, fmt.Sprintf("PHOTON TORPEDOES: %d", q.Player.Torpedoes))
	game.EmitStr(49, 10, fmt.Sprintf("KLINGONS:         %d", q.Game.GetRemainingKlingons()))
}

func newLocation(ox int, oy int, direction int) (int, int) {
	x, y := ox, oy
	switch direction {
	case Dir1:
		x--
		y++
	case Dir2:
		y++
	case Dir3:
		x++
		y++
	case Dir4:
		x--
	case Dir5:
		// no movement
	case Dir6:
		x++
	case Dir7:
		x--
		y--
	case Dir8:
		y--
	case Dir9:
		x++
		y--
	}
	return x, y
}

// MoveObject moves the object in the direction specified
func (q *Quadrant) MoveObject(m MoveableObject, direction int) {
	ox, oy := m.Location()
	x, y := newLocation(ox, oy, direction)

	if x < 0 {
		x = 0
	} else if x > 9 {
		x = 9
	}

	if y < 0 {
		y = 0
	} else if y > 9 {
		y = 9
	}

	if q.Objects[x][y] == nil {
		q.Objects[x][y] = m
		q.Objects[ox][oy] = nil
		m.Move(x, y)
	}
}
