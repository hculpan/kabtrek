package quadrant

// Player interface to test if object is player
type Player interface {
	MoveableObject
	IsPlayer() bool
}

// MoveableObject is the interface for all moveable objects
type MoveableObject interface {
	Object
	Move(x int, y int)
}

// Object interface to generize all objects in sector
type Object interface {
	Name() string
	Location() (int, int)
	TakeDamage(damage int)
	GetShields() int
}
