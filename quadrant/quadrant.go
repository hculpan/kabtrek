package quadrant

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/hculpan/kabtrek/util"
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
	Navigation
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
	Stardate float32
}

// Quadrant : All the information related to a single Quadrant
type Quadrant struct {
	X                        int
	Y                        int
	Objects                  [10][10]Object
	Player                   *Enterprise
	Stardate                 float32
	StartingNumberOfKlingons int
	NumberOfKlingons         int

	UIState       int
	AwaitingInput bool
	CurrentInput  string

	Messages []Message

	// Private variables
	blinkRed  int
	torpedoes [10][10]*Torpedo
}

// NewQuadrant creates a new quadrant, populated with items
func NewQuadrant(x int, y int, numKlingons int, numStars int, numBases int) *Quadrant {
	result := &Quadrant{
		X:                        x,
		Y:                        y,
		Objects:                  [10][10]Object{},
		Player:                   nil,
		Stardate:                 3100.1,
		NumberOfKlingons:         numKlingons,
		StartingNumberOfKlingons: numKlingons,
		UIState:                  0,
		AwaitingInput:            false,
		CurrentInput:             "",
		blinkRed:                 0,
		torpedoes:                [10][10]*Torpedo{},
	}

	result.blinkRed = 0
	result.UIState = Normal
	result.AwaitingInput = false
	result.CurrentInput = ""

	result.NumberOfKlingons = numKlingons
	klingonsToPlace := numKlingons
	for klingonsToPlace > 0 {
		xloc := util.RandomInt(10)
		yloc := util.RandomInt(10)
		if result.Objects[xloc][yloc] == nil {
			result.Objects[xloc][yloc] = NewKlingon(xloc, yloc)
			klingonsToPlace--
		}
	}

	starsToPlace := numStars
	for starsToPlace > 0 {
		xloc := util.RandomInt(10)
		yloc := util.RandomInt(10)
		if result.Objects[xloc][yloc] == nil {
			result.Objects[xloc][yloc] = &Star{X: xloc, Y: yloc}
			starsToPlace--
		}
	}

	basesToPlace := numBases
	for basesToPlace > 0 {
		xloc := util.RandomInt(10)
		yloc := util.RandomInt(10)
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
func (q *Quadrant) UpdateState(scr tcell.Screen, newState int) {
	q.UIState = newState
	if newState != Normal && newState != Weapons {
		q.AwaitingInput = true
	} else {
		q.AwaitingInput = false
	}
	q.CurrentInput = ""
	q.Draw(scr)
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

func (q *Quadrant) damageObjectAt(x int, y int, damage int, descriptor string) {
	if q.Objects[x][y] == nil {
		return
	}

	q.Objects[x][y].TakeDamage(TorpedoDamage)
	q.AddMessage(fmt.Sprintf("%s at %d, %d took %d damage from a %s", q.Objects[x][y].Name(), x, y, damage, descriptor))

	if q.Objects[x][y].GetShields() <= 0 {
		q.AddMessage(fmt.Sprintf("%s at %d, %d destroyed!", q.Objects[x][y].Name(), x, y))
		switch q.Objects[x][y].(type) {
		case *Klingon:
			q.NumberOfKlingons--
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
	if util.CheckPercent(66) {
		if util.CheckPercent(25) {
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
			q.MoveObject(o, util.RandomInt(9)+1)
		}
	}
}

// AddMessage adds a message to the messages display
// Will be removed in 5 turns (stardate + 0.5)
func (q *Quadrant) AddMessage(t string) {
	q.Messages = append(q.Messages, Message{Text: t, Stardate: q.Stardate})
}

// Update processes the next turn for the quadrant
func (q *Quadrant) Update() {
	q.Stardate += 0.1

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			switch q.Objects[x][y].(type) {
			case *Klingon:
				q.klingonAction(q.Objects[x][y].(*Klingon))
			}
		}
	}

	if q.playerDockedAtBase() {
		q.Player.Energy += EnterpriseMaxEnergy * 0.25
		if q.Player.Energy > EnterpriseMaxEnergy {
			q.Player.Energy = EnterpriseMaxEnergy
		}
		q.Player.Torpedoes += EnterpriseMaxTorpedoes * 0.25
		if q.Player.Torpedoes > EnterpriseMaxTorpedoes {
			q.Player.Torpedoes = EnterpriseMaxTorpedoes
		}
	}
}

// Draw draw's the quadrant
func (q *Quadrant) Draw(scr tcell.Screen) {
	scr.Sync()
	scr.Clear()
	q.DisplayQuadrant(scr)
	q.DisplayStatus(scr)
	q.displayState(scr)
	q.displayMessages(scr)
	scr.Show()
}

// DisplayQuadrant draws the Quadrant map
func (q *Quadrant) DisplayQuadrant(scr tcell.Screen) {
	QuadrantStr := fmt.Sprintf("Quadrant : %d, %d", q.X+1, q.Y+1)
	util.EmitStr(scr, 22-(len(QuadrantStr)/2), 0, tcell.StyleDefault, QuadrantStr)
	util.EmitStr(scr, 3, 1, tcell.StyleDefault, "=---=---=---=---=---=---=---=---=---=---")
	for i := 0; i < 10; i++ {
		q.displayQuadrantLine(scr, i)
	}
	util.EmitStr(scr, 3, 12, tcell.StyleDefault, "=-1-=-2-=-3-=-4-=-5-=-6-=-7-=-8-=-9-=-10")
}

func (q *Quadrant) displayMessages(scr tcell.Screen) {
	for i, t := range q.Messages {
		util.EmitStr(scr, 1, 16+i, tcell.StyleDefault, fmt.Sprintf("Stardate %.1f: %s", t.Stardate, t.Text))
	}

	// Remove older messages
	for len(q.Messages) > 0 {
		if q.Messages[0].Stardate < q.Stardate-0.5 {
			q.Messages = q.Messages[1:]
		} else {
			break
		}
	}
}

func (q *Quadrant) displaySector(scr tcell.Screen, x int, y int) {
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
		util.EmitStr(scr, x*4+4, y+2, tcell.StyleDefault, objStr)
	} else if q.torpedoes[x][y] != nil {
		util.EmitStr(scr, x*4+4, y+2, tcell.StyleDefault, " @ ")
	}

}

func (q *Quadrant) displayQuadrantLine(scr tcell.Screen, row int) {
	util.EmitStr(scr, 0, row+2, tcell.StyleDefault, fmt.Sprintf("%2d|                                        |", row+1))

	for i := 0; i < 10; i++ {
		q.displaySector(scr, i, row)
	}
}

// HandleKeyForState handles the current key for the various ui states
func (q *Quadrant) HandleKeyForState(scr tcell.Screen, key tcell.EventKey) {
	switch q.UIState {
	case Weapons:
		switch int(key.Rune()) {
		case 'p', 'P':
			q.UpdateState(scr, WeaponsPhasers)
			q.Draw(scr)
		case 't', 'T':
			q.UpdateState(scr, WeaponsTorpedoes)
			q.Draw(scr)
		}
	}
}

// AcceptInput accepts whatever the player has typed for input
func (q *Quadrant) AcceptInput(scr tcell.Screen) {
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
	q.UpdateState(scr, Normal)
	q.Draw(scr)
}

func (q *Quadrant) displayState(scr tcell.Screen) {
	switch q.UIState {
	case Normal:
		util.EmitStr(scr, 1, 14, tcell.StyleDefault, "(N)avigation   (W)eapons   (S)hields  S(e)nsors  Ship's (C)omputer")
	case Shields:
		util.EmitStr(scr, 1, 14, tcell.StyleDefault, "Set energy for shields: ")
		q.displayInput(scr, 25, 14)
	case Weapons:
		util.EmitStr(scr, 1, 14, tcell.StyleDefault, "(P)hasers or Photon (T)orpedoes")
	case WeaponsTorpedoes:
		util.EmitStr(scr, 1, 14, tcell.StyleDefault, "Direction:")
		q.displayInput(scr, 12, 14)
	}
}

func (q *Quadrant) displayInput(scr tcell.Screen, col int, row int) {
	util.EmitStr(scr, col, row, tcell.StyleDefault, q.CurrentInput)
	if q.blinkRed%2 == 0 {
		util.EmitStr(scr, col+len(q.CurrentInput), row, tcell.StyleDefault, "_")
	}
}

// DisplayStatus draws the status of the quadrant (the stuff to the right of the map)
func (q *Quadrant) DisplayStatus(scr tcell.Screen) {
	q.blinkRed++
	util.EmitStr(scr, 46, 3, tcell.StyleDefault, fmt.Sprintf("STARDATE:         %.1f", q.Stardate))
	util.EmitStr(scr, 46, 4, tcell.StyleDefault, fmt.Sprintf("SECTOR:           %d,%d", q.Player.X, q.Player.Y))

	if q.playerDockedAtBase() {
		util.EmitStr(scr, 46, 5, tcell.StyleDefault, "CONDITION:        DOCKED")
	} else if q.NumberOfKlingons > 0 && q.blinkRed%2 == 1 {
		util.EmitStr(scr, 46, 5, tcell.StyleDefault, "CONDITION:        RED")
	} else if q.NumberOfKlingons > 0 {
		util.EmitStr(scr, 46, 5, tcell.StyleDefault, "CONDITION: ")
	} else {
		util.EmitStr(scr, 46, 5, tcell.StyleDefault, "CONDITION:        GREEN")
	}

	util.EmitStr(scr, 46, 6, tcell.StyleDefault, fmt.Sprintf("SHIELDS:          %d", q.Player.Shields))
	util.EmitStr(scr, 46, 7, tcell.StyleDefault, fmt.Sprintf("ENERGY:           %d", q.Player.Energy))
	util.EmitStr(scr, 46, 8, tcell.StyleDefault, fmt.Sprintf("PHOTON TORPEDOES: %d", q.Player.Torpedoes))
	util.EmitStr(scr, 46, 10, tcell.StyleDefault, fmt.Sprintf("KLINGONS:         %d", q.NumberOfKlingons))
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
