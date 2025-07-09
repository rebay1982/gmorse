package fft

import (
	"math"
	"testing"

	"github.com/rebay1982/gmorse/internal/test"
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
		// Review this case, some complex numbers don't have the correct sign after FFT computation.
		{
			name:    "recursive_twiddle_exercice_input",
			samples: []complex128{0, 1, 0, 0, 0, 0, 0, 0},
			expected: []complex128{
				complex(1, 0),
				complex(math.Cos(-math.Pi/4*1), math.Sin(-math.Pi/4*1)),
				complex(math.Cos(-math.Pi/4*2), math.Sin(-math.Pi/4*2)),
				complex(math.Cos(-math.Pi/4*3), math.Sin(-math.Pi/4*3)),
				complex(math.Cos(-math.Pi/4*4), math.Sin(-math.Pi/4*4)),
				complex(math.Cos(-math.Pi/4*5), math.Sin(-math.Pi/4*5)),
				complex(math.Cos(-math.Pi/4*6), math.Sin(-math.Pi/4*6)),
				complex(math.Cos(-math.Pi/4*7), math.Sin(-math.Pi/4*7)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fft := RecursiveFFT(tc.samples)
			test.ValidateFFT(t, tc.expected, fft)
		})
	}
}
