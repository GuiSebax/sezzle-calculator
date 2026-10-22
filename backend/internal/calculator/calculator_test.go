package calculator

import (
	"errors"
	"math"
	"testing"
)

type binaryCase struct {
	name    string
	a, b    float64
	want    float64
	wantErr error
}

func runBinary(t *testing.T, fn func(a, b float64) (float64, error), tests []binaryCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fn(tt.a, tt.b)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func assertResult(t *testing.T, got float64, err error, want float64, wantErr error) {
	t.Helper()
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if wantErr != nil {
		return
	}
	if !almostEqual(got, want) {
		t.Errorf("result = %v, want %v", got, want)
	}
}

// almostEqual compares floats with a relative tolerance to absorb rounding error.
func almostEqual(a, b float64) bool {
	const eps = 1e-9
	if a == b {
		return true
	}
	return math.Abs(a-b) <= eps*math.Max(math.Abs(a), math.Abs(b))
}

func TestAdd(t *testing.T) {
	runBinary(t, Add, []binaryCase{
		{name: "positive numbers", a: 5, b: 3, want: 8},
		{name: "negative numbers", a: -5, b: -3, want: -8},
		{name: "mixed signs", a: -5, b: 3, want: -2},
		{name: "zeros", a: 0, b: 0, want: 0},
		{name: "decimals", a: 0.1, b: 0.2, want: 0.3},
		{name: "overflow", a: math.MaxFloat64, b: math.MaxFloat64, wantErr: ErrNonFinite},
	})
}

func TestSubtract(t *testing.T) {
	runBinary(t, Subtract, []binaryCase{
		{name: "positive result", a: 5, b: 3, want: 2},
		{name: "negative result", a: 3, b: 5, want: -2},
		{name: "subtract negative", a: 5, b: -3, want: 8},
		{name: "decimals", a: 0.3, b: 0.1, want: 0.2},
		{name: "overflow", a: -math.MaxFloat64, b: math.MaxFloat64, wantErr: ErrNonFinite},
	})
}

func TestMultiply(t *testing.T) {
	runBinary(t, Multiply, []binaryCase{
		{name: "positive numbers", a: 5, b: 3, want: 15},
		{name: "by zero", a: 5, b: 0, want: 0},
		{name: "negative and positive", a: -5, b: 3, want: -15},
		{name: "two negatives", a: -5, b: -3, want: 15},
		{name: "decimals", a: 1.5, b: 2.5, want: 3.75},
		{name: "overflow", a: math.MaxFloat64, b: 2, wantErr: ErrNonFinite},
	})
}

func TestDivide(t *testing.T) {
	runBinary(t, Divide, []binaryCase{
		{name: "exact", a: 6, b: 3, want: 2},
		{name: "fractional", a: 1, b: 3, want: 1.0 / 3},
		{name: "negative divisor", a: 6, b: -3, want: -2},
		{name: "zero dividend", a: 0, b: 5, want: 0},
		{name: "by zero", a: 5, b: 0, wantErr: ErrDivisionByZero},
		{name: "zero by zero", a: 0, b: 0, wantErr: ErrDivisionByZero},
		{name: "by negative zero", a: 5, b: math.Copysign(0, -1), wantErr: ErrDivisionByZero},
		{name: "overflow", a: math.MaxFloat64, b: 0.5, wantErr: ErrNonFinite},
	})
}

func TestPower(t *testing.T) {
	runBinary(t, Power, []binaryCase{
		{name: "integer exponent", a: 2, b: 10, want: 1024},
		{name: "zero exponent", a: 5, b: 0, want: 1},
		{name: "zero to zero", a: 0, b: 0, want: 1},
		{name: "negative exponent", a: 2, b: -1, want: 0.5},
		{name: "fractional exponent", a: 9, b: 0.5, want: 3},
		{name: "negative base odd exponent", a: -2, b: 3, want: -8},
		{name: "zero to negative exponent", a: 0, b: -1, wantErr: ErrNonFinite},
		{name: "negative base fractional exponent", a: -8, b: 1.0 / 3, wantErr: ErrNonFinite},
		{name: "overflow", a: 10, b: 400, wantErr: ErrNonFinite},
	})
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		want    float64
		wantErr error
	}{
		{name: "perfect square", a: 16, want: 4},
		{name: "irrational", a: 2, want: math.Sqrt2},
		{name: "zero", a: 0, want: 0},
		{name: "fraction", a: 0.25, want: 0.5},
		{name: "negative", a: -1, wantErr: ErrNegativeSqrt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sqrt(tt.a)
			assertResult(t, got, err, tt.want, tt.wantErr)
		})
	}
}

func TestPercentage(t *testing.T) {
	runBinary(t, Percentage, []binaryCase{
		{name: "ten percent of fifty", a: 10, b: 50, want: 5},
		{name: "zero percent", a: 0, b: 50, want: 0},
		{name: "hundred percent", a: 100, b: 42, want: 42},
		{name: "over hundred percent", a: 150, b: 20, want: 30},
		{name: "negative percent", a: -10, b: 50, want: -5},
		{name: "fractional percent", a: 12.5, b: 80, want: 10},
		{name: "overflow", a: math.MaxFloat64, b: math.MaxFloat64, wantErr: ErrNonFinite},
	})
}
