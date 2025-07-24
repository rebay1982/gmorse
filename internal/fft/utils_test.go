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
			got := NormalizePCM16Samples(tc.samples)

			validateNormalizedPCMData(t, tc.expected, got)
		})
	}
}

func Test_IsPowerOfTwo(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected bool
	}{
		{
			name:     "0",
			input:    0,
			expected: false,
		},
		{
			name:     "1",
			input:    1,
			expected: true,
		},
		{
			name:     "2",
			input:    2,
			expected: true,
		},
		{
			name:     "3",
			input:    3,
			expected: false,
		},
		{
			name:     "4",
			input:    4,
			expected: true,
		},
		{
			name:     "8",
			input:    8,
			expected: true,
		},
		{
			name:     "15",
			input:    15,
			expected: false,
		},
		{
			name:     "16",
			input:    16,
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := isPowerOfTwo(tc.input)

			if tc.expected != got {
				t.Errorf("For input %d: Expected %t, got %t\n", tc.input, tc.expected, got)
			}
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
