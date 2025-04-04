package calculator

import (
	"fmt"
	"math"
)

func Add(a, b float64) (float64, error) {
	if b > 0 {
		if a > math.MaxFloat64-b {
			return 0, fmt.Errorf("got overflow")
		}
	} else {
		if a < -math.MaxFloat64-b {
			return 0, fmt.Errorf("got overflow")
		}
	}
	return a + b, nil
}

func Sub(a, b float64) (float64, error) {
	return Add(a, -b)
}

func Mul(a, b float64) (float64, error) {
	if math.Abs(a) > math.MaxFloat64/math.Abs(b) {
		return 0, fmt.Errorf("got overflow")
	}
	return a * b, nil
}

func Div(a, b float64) (float64, error) {
	return Mul(a, 1/b)
}

func Pow(a, b float64) (float64, error) {
	if b > 1 && a > 1 && b > Log(a, math.MaxFloat64) {
		return 0, fmt.Errorf("got overflow")
	}
	return math.Pow(a, b), nil
}

func Log(base, x float64) float64 {
	return math.Log(x) / math.Log(base)
}