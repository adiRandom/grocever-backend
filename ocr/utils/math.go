package utils

import "math"

func TruncateFloat(f float64, n int) float64 {
	shift := math.Pow(10, float64(n))
	return math.Trunc((f+0.5/shift)*shift) / shift
}
