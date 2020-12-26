package galaxy

import (
	"github.com/hculpan/kabtrek/quadrant"
	"github.com/hculpan/kabtrek/util"
)

// Galaxy contains info for, well, the galaxy
type Galaxy struct {
	Stardate                  float64
	Quadrants                 [8][6]quadrant.Quadrant
	StartingNumberOfKlingons  int
	NumberOfKlingons          int
	StartingNumberOfStarbases int
	NumberOfStarbases         int
	Player                    *quadrant.Enterprise
}

// NewGalaxy create a whole new galaxy
func NewGalaxy(numKlingons, numStarbases int) *Galaxy {
	result := &Galaxy{
		Stardate:                  3700.1,
		Quadrants:                 [8][6]quadrant.Quadrant{},
		StartingNumberOfKlingons:  numKlingons,
		NumberOfKlingons:          numKlingons,
		StartingNumberOfStarbases: numStarbases,
		NumberOfStarbases:         numStarbases,
	}

	quadsGened := [8][6]bool{}
	for x := 0; x < 8; x++ {
		for y := 0; y < 6; y++ {
			quadsGened[x][y] = false
		}
	}

	remainingKlingons := numKlingons
	remainingStarbases := numStarbases
	for remainingKlingons > 0 {
		// Pick number of klingons to place
		numKlingonsInQuadrant := selectKlingonsForQuadrant()
		if numKlingonsInQuadrant > remainingKlingons {
			numKlingonsInQuadrant = remainingKlingons
		}
		remainingKlingons -= numKlingonsInQuadrant

		for {
			x := util.RandomInt(8)
			y := util.RandomInt(6)
			if !quadsGened[x][y] && remainingStarbases > 0 && util.CheckPercent(10) {
				result.Quadrants[x][y] = *quadrant.NewQuadrant(x, y, numKlingonsInQuadrant, util.RandomInt(7), 1)
				remainingStarbases--
				quadsGened[x][y] = true
				break
			} else if !quadsGened[x][y] {
				result.Quadrants[x][y] = *quadrant.NewQuadrant(x, y, numKlingonsInQuadrant, util.RandomInt(7), 0)
				quadsGened[x][y] = true
				break
			}
		}

	}

	for remainingStarbases > 0 {
		x := util.RandomInt(8)
		y := util.RandomInt(6)
		if !quadsGened[x][y] {
			result.Quadrants[x][y] = *quadrant.NewQuadrant(x, y, 0, util.RandomInt(7), 1)
			remainingStarbases--
			quadsGened[x][y] = true
		}
	}

	for x := 0; x < 8; x++ {
		for y := 0; y < 6; y++ {
			if !quadsGened[x][y] {
				result.Quadrants[x][y] = *quadrant.NewQuadrant(x, y, 0, util.RandomInt(7), 0)
			}
		}
	}

	return result
}

func selectKlingonsForQuadrant() int {
	n := util.GetPercent()
	switch {
	case n < 55:
		return 0
	case n < 75:
		return 1
	case n < 85:
		return 2
	case n < 92:
		return 3
	case n < 97:
		return 4
	default:
		return 5
	}
}
