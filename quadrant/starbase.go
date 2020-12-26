package quadrant

// Starbase represents planets in the sector
type Starbase struct {
	X       int
	Y       int
	Shields int
}

// NewStarbase returns a new instance of starbase
func NewStarbase(x int, y int) *Starbase {
	return &Starbase{
		X:       x,
		Y:       y,
		Shields: 10000}
}

// Location returns the location of the planet in the quadrant
func (s Starbase) Location() (int, int) {
	return s.X, s.Y
}

// TakeDamage does damage to klingon
func (s *Starbase) TakeDamage(damage int) {
	s.Shields -= damage
}

// GetShields returns the object's shield strength
func (s Starbase) GetShields() int {
	return s.Shields
}

// Name returns the display-friendly name
func (s Starbase) Name() string {
	return "Starbase"
}
