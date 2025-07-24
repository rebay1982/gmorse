package windowing

import (
	"math"
	"testing"

	"github.com/rebay1982/gmorse/internal/test"
)

func Test_Hann(t *testing.T) {
	testCases := []struct {
		name     string
		samples  []float64
		expected []float64
	}{
		{
			name:    "constant_even_samples",
			samples: []float64{1, 1, 1, 1, 1, 1, 1, 1},
			expected: []float64{
				0.5 * (1 - math.Cos(2*math.Pi*float64(0)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(1)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(2)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(3)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(4)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(5)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(6)/float64(7))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(7)/float64(7))),
			},
		},
		{
			// Should have a peek at 1
			name:    "constant_odd_samples",
			samples: []float64{1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: []float64{
				0.5 * (1 - math.Cos(2*math.Pi*float64(0)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(1)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(2)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(3)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(4)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(5)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(6)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(7)/float64(8))),
				0.5 * (1 - math.Cos(2*math.Pi*float64(8)/float64(8))),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := make([]float64, len(tc.samples))
			copy(got, tc.samples)

			Hann(got)

			validateHanningWindowing(t, tc.expected, got)
		})
	}
}

func validateHanningWindowing(t *testing.T, expected, samples []float64) {
	expectedLen := len(expected)
	samplesLen := len(samples)

	if expectedLen != samplesLen {
		t.Errorf("Expected length %d, got %d", expectedLen, samplesLen)
	}

	for i, s := range samples {
		if !test.Approximately(s, expected[i]) {
			t.Errorf("At frequency bin %d: Expected %.7f, got %.7f\n", i, expected[i], s)
		}
	}
}
