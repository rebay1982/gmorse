package fft

import (
	"testing"
)

func Test_RecursiveFFT(t *testing.T) {
	testCases := []struct {
		name     string
		samples  []complex128
		expected []complex128
	}{
		{
			name:     "recursive_impulse_input",
			samples:  []complex128{1, 0, 0, 0, 0, 0, 0, 0},
			expected: []complex128{1, 1, 1, 1, 1, 1, 1, 1},
		},
		{
			name:     "recursive_dc_input",
			samples:  []complex128{1, 1, 1, 1, 1, 1, 1, 1},
			expected: []complex128{8, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "recursive_single_freq_input",
			samples:  []complex128{1, -1, 1, -1, 1, -1, 1, -1},
			expected: []complex128{0, 0, 0, 0, 8, 0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			fft := RecursiveFFT(tc.samples)

			validateRecursiveFFT(t, tc.expected, fft)
		})
	}
}

func validateRecursiveFFT(t *testing.T, expected, fft []complex128) {
	expectedLen := len(expected)
	fftLen := len(fft)

	if expectedLen != fftLen {
		t.Errorf("Expected length %d, got %d", expectedLen, fftLen)
	}

	for k, x := range fft {
		if x != expected[k] {
			t.Errorf("At frequency bin %d: Expected %.2f, got %.2f\n", k, expected[k], x)
		}
	}
}
