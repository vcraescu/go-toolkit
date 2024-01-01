package math

import "math"

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))

	num *= output
	num = float64(int(num + math.Copysign(0.5, num)))

	return num / output
}
