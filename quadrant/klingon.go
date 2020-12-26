package quadrant

// Klingon for the Enterprise to blow up
type Klingon struct {
	X         int
	Y         int
	Shields   int
	Torpedoes int
}

// NewKlingon creates a new Klingon
func NewKlingon(x int, y int) *Klingon {
	return &Klingon{X: x, Y: y, Shields: 1000, Torpedoes: 10}
}

// Move the klingon
func (k *Klingon) Move(x int, y int) {
	k.X = x
	k.Y = y
}

// TakeDamage does damage to klingon
func (k *Klingon) TakeDamage(damage int) {
	k.Shields -= damage
}

// Location returns the location of the Enterprise
func (k Klingon) Location() (int, int) {
	return k.X, k.Y
}

// GetShields returns the object's shield strength
func (k Klingon) GetShields() int {
	return k.Shields
}

// Name returns the display-friendly name
func (k Klingon) Name() string {
	return "Klingon"
}
