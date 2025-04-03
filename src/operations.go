package main

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
