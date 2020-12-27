package game

import (
	"math/rand"
	"time"
)

// Rnd is global random
var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomInt returns a random int between 0 and n-1, inclusive
func RandomInt(n int) int {
	return rnd.Intn(n)
}

// CheckPercent returns true/false based on random number weighted by percentage
func CheckPercent(percentage int) bool {
	return rnd.Intn(100) < percentage
}

// GetPercent returns a number between 0 and 99
func GetPercent() int {
	return rnd.Intn(100)
}
