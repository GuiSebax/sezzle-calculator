// Package calculator implements basic arithmetic operations as pure functions.
//
// Every function returns an error instead of producing NaN or ±Inf, so callers
// never have to inspect the numeric result to know whether it is usable.
package calculator

import (
	"errors"
	"math"
)

// Sentinel errors returned by the operations. Callers should match them with errors.Is.
var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrNegativeSqrt   = errors.New("square root of a negative number")
	ErrNonFinite      = errors.New("result is not a finite number")
)

// Add returns a + b.
func Add(a, b float64) (float64, error) {
	return finite(a + b)
}

// Subtract returns a - b.
func Subtract(a, b float64) (float64, error) {
	return finite(a - b)
}

// Multiply returns a * b.
func Multiply(a, b float64) (float64, error) {
	return finite(a * b)
}

// Divide returns a / b, or ErrDivisionByZero when b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return finite(a / b)
}

// Power returns a raised to the power b.
// Results without a real value (e.g. a negative base with a fractional
// exponent) or that overflow return ErrNonFinite.
func Power(a, b float64) (float64, error) {
	return finite(math.Pow(a, b))
}

// Sqrt returns the square root of a, or ErrNegativeSqrt when a is negative.
func Sqrt(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSqrt
	}
	return finite(math.Sqrt(a))
}

// Percentage returns a percent of b (e.g. Percentage(10, 50) == 5).
func Percentage(a, b float64) (float64, error) {
	// Multiplying before dividing keeps integer inputs exact (10*50/100 == 5).
	return finite(a * b / 100)
}

// finite returns ErrNonFinite if v is NaN or ±Inf.
func finite(v float64) (float64, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, ErrNonFinite
	}
	return v, nil
}
