package goertzel

import (
	"github.com/rebay1982/gmorse/internal/fft"
	"github.com/rebay1982/gmorse/internal/test"
	"testing"
)

func Test_Goertzel(t *testing.T) {
	testCases := []struct {
		name       string
		samples    []float64
		sampleRate float64
		targetFreq float64
		expected   float64
	}{
		{
			// Absence of signal
			name:       "no_signal",
			samples:    []float64{0, 0, 0, 0, 0, 0, 0, 0},
			sampleRate: 8.0,
			targetFreq: 0.0,
			expected:   0.0,
		},
		{
			// Constant signal (DC)
			name:       "dc",
			samples:    []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			sampleRate: 8.0,
			targetFreq: 2.0,
			expected:   0.0,
		},
		{
			// 1Hz sine signal.
			name:       "sine",
			samples:    []float64{0.0, 0.707, 1.0, 0.707, 0.0, -0.707, -1.0, -0.707},
			sampleRate: 8.0,
			targetFreq: 1.0,
			expected:   3.999697,
		},
		{
			// 3Hz sine signal.
			name:       "sine_off_frequency",
			samples:    []float64{0.0, 0.707, -1.0, -0.707, 0.0, -0.707, 1.0, -0.707},
			sampleRate: 8.0,
			targetFreq: 2.0,
			expected:   1.4140,
		},
		{
			// 1Hz and 3Hz sine signal.
			name:       "sine_two_diff_frequencies",
			samples:    []float64{0.0, 1.414, 0.0, 0.0, 0.0, -1.414, 0, -1.414},
			sampleRate: 8.0,
			targetFreq: 3.0,
			expected:   3.161800,
		},
		{
			// Constant signal (DC)
			name:       "noise",
			samples:    []float64{0.53, 0.13, 0.89, 0.20, 0.26, 0.78, 0.66, 0.01, 0.43},
			sampleRate: 8.0,
			targetFreq: 2.0,
			expected:   0.752888,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			output := Goertzel(tc.sampleRate, tc.targetFreq, tc.samples)
			got := fft.ComputeMagnitude(output)

			if !test.Approximately(tc.expected, got) {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}
