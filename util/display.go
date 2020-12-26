package util

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Rnd is global random
var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// EmitStr will print a string to the given screen
func EmitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

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
