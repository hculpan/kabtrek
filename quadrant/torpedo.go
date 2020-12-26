package quadrant

// Torpedo is the struct representing a photon torpedo
type Torpedo struct {
	X         int
	Y         int
	Direction int
}

// Move the torpedo
func (t *Torpedo) Move(x int, y int) {
	t.X = x
	t.Y = y
}

// Location returns the location of the torpedo
func (t Torpedo) Location() (int, int) {
	return t.X, t.Y
}
