package game

// Constants related to player's ship
const (
	EnterpriseMaxEnergy    = 5000
	EnterpriseMaxTorpedoes = 20
)

// Constants for game UI state
const (
	Quadrant = iota
	GalaxyMap
	LongRangeSensors
	Quitting
)

// QuadrantSummary gives summary of quadrant
// that is used to display galaxy map
type QuadrantSummary struct {
	X         int
	Y         int
	Klingons  int
	Starbases int
	Stars     int
	IsActive  bool
	Scanned   bool
}

// Game is the global object with all the overall game state
type Game interface {
	GetStartingKlingons() int
	GetRemainingKlingons() int
	KlingonDestroyed()
	GetStartingStarbases() int
	GetRemainingStarbases() int
	StarbaseDestroyed()

	GetStardate() float64

	GetQuadrantSummary(x, y int) *QuadrantSummary

	SetGameState(state int)

	NavigateTo(x, y int)

	Draw()
}
