package game

import "math"

// Distance calculates the distance between two points
func Distance(x1, y1, x2, y2 int) float64 {
	first := math.Pow(float64(x2-x1), 2)
	second := math.Pow(float64(y2-y1), 2)
	return math.Sqrt(first + second)
}
