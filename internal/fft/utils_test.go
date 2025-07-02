package fft

import (
	"math"
	"testing"
)

func validateFFT(t *testing.T, expected, fft []complex128) {
	expectedLen := len(expected)
	fftLen := len(fft)

	if expectedLen != fftLen {
		t.Errorf("Expected length %d, got %d", expectedLen, fftLen)
	}

	for k, x := range fft {
		if !approxEqualComplex(x, expected[k], 0.0001) {
			t.Errorf("At frequency bin %d: Expected %.2f, got %.2f\n", k, expected[k], x)
		}
	}
}

func approxEqualComplex(x, y complex128, tolerance float64) bool {
	return math.Abs(real(x)-real(y)) < tolerance && math.Abs(imag(x)-imag(y)) < tolerance
}
