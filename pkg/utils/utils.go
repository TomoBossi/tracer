package utils

import "math"

func Abs(n int) int {
	if n < 0 {
		return int(-n)
	}
	return int(n)
}

func L2(x1, y1, x2, y2 int) float64 {
	dx, dy := x1-x2, y1-y2
	return math.Sqrt(float64(dx*dx + dy*dy))
}
