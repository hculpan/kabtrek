package quadrant

import "math"

// Star represents planets in the sector
type Star struct {
	X int
	Y int
}

// Location returns the location of the star in the quadrant
func (s Star) Location() (int, int) {
	return s.X, s.Y
}

// TakeDamage does no damage to star
func (s *Star) TakeDamage(damage int) {
	// stars cannot be damaged
}

// GetShields returns the object's shield strength
func (s Star) GetShields() int {
	return math.MaxInt64
}

// Name returns the display-friendly name
func (s Star) Name() string {
	return "Star"
}
