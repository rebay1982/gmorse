package fft

import (
	"testing"

	"github.com/rebay1982/gmorse/internal/test"
)

func Test_Normalize16BitPCM(t *testing.T) {
	testCases := []struct {
		name     string
		samples  []int16
		expected []float64
	}{
		{
			name:     "max_positive_16bit_pcm",
			samples:  []int16{int16(32767)},
			expected: []float64{0.99997},
		},
		{
			name:     "min_negative_16bit_pcm",
			samples:  []int16{int16(-32768)},
			expected: []float64{-1.0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizePCM16(tc.samples)

			validateNormalizedPCMData(t, tc.expected, got)
		})
	}
}

func validateNormalizedPCMData(t *testing.T, expected, normalized []float64) {
	expectedLen := len(expected)
	normalizedLen := len(normalized)

	if expectedLen != normalizedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, normalizedLen)
	}

	for i, n := range normalized {
		if !test.Approximately(n, expected[i]) {
			t.Errorf("At frequency bin %d: Expected %.7f, got %.7f\n", i, expected[i], n)
		}
	}
}
