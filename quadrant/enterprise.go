package quadrant

// EnterpriseMaxEnergy is the maximum energy for the Enterprise
const (
	EnterpriseMaxEnergy    = 5000
	EnterpriseMaxTorpedoes = 20
)

// Enterprise : Information relating to the player ship
type Enterprise struct {
	X         int
	Y         int
	QuadrantX int
	QuadrantY int
	Energy    int
	Shields   int
	Torpedoes int
}

// NewEnterprise creates a new Enterprise
func NewEnterprise(xloc int, yloc int) *Enterprise {
	return &Enterprise{
		X:         xloc,
		Y:         yloc,
		Energy:    EnterpriseMaxEnergy,
		Shields:   0,
		Torpedoes: EnterpriseMaxTorpedoes,
	}
}

// Move the enterprise
func (e *Enterprise) Move(x int, y int) {
	if e.X != x || e.Y != y {
		e.Energy -= EnergyToMove + ((e.Shields / 1000) * EnergyToMove)
	}

	e.X = x
	e.Y = y
}

// TakeDamage reduces damage to Enterprise
func (e *Enterprise) TakeDamage(damage int) {
	if damage <= e.Shields {
		e.Shields -= damage
	} else {
		damage -= e.Shields
		e.Shields = 0
		e.Energy -= damage * 2
	}
}

// Location returns the location of the Enterprise
func (e Enterprise) Location() (int, int) {
	return e.X, e.Y
}

// IsPlayer indicates if this is the player object or not
func (e *Enterprise) IsPlayer() bool {
	return true
}

// GetShields returns the object's shield strength
func (e Enterprise) GetShields() int {
	if e.Shields > 0 {
		return e.Shields
	} else {
		return e.Energy
	}
}

// Name returns the display-friendly name
func (e Enterprise) Name() string {
	return "Enterprise"
}
