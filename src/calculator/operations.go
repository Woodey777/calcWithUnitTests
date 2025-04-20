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
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
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

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, fmt.Errorf("square root of negative number")
	}
	return math.Sqrt(x), nil
}

func Ln(x float64) (float64, error) {
	if x <= 0 {
		return 0, fmt.Errorf("natural logarithm of non-positive number")
	}
	return math.Log(x), nil
}

func Sin(x float64) float64 {
	return math.Sin(x)
}

func Cos(x float64) float64 {
	return math.Cos(x)
}

func Tg(x float64) (float64, error) {
	if _, frac := math.Modf(x / (math.Pi/2)); frac == 0 {
		return 0, fmt.Errorf("tangent of pi/2 * k")
	}
	return math.Tan(x), nil
}

func Cot(x float64) (float64, error) {
	if _, frac := math.Modf(x / math.Pi); frac == 0  {
		return 0, fmt.Errorf("cotangent of pi * k")
	}
	return 1 / math.Tan(x), nil
}

func Exp(x float64) (float64, error) {
	if x > math.Log(math.MaxFloat64) {
		return 0, fmt.Errorf("got overflow")
	}
	return math.Exp(x), nil
}

