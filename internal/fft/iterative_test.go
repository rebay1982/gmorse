package fft

import (
	"math"
	"testing"

	"github.com/rebay1982/gmorse/internal/test"
)

func Test_reverseBits(t *testing.T) {
	testCases := []struct {
		name     string
		in       uint
		maskSize uint // nb bits
		expected uint
	}{
		{
			name:     "0",
			in:       0,
			maskSize: 32,
			expected: 0,
		},
		{
			name:     "1",
			in:       1,
			maskSize: 3,
			expected: 4,
		},
		{
			name:     "3",
			in:       3,
			maskSize: 3,
			expected: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := reverseBits(tc.in, tc.maskSize)

			if tc.expected != got {
				t.Errorf("Expected %d for input %d, got %d", tc.expected, tc.in, got)
			}
		})
	}
}

func Test_BitReverseIndexPrecomputation(t *testing.T) {
	expected := []int{0, 4, 2, 6, 1, 5, 3, 7}

	got := preComputeBitReverseIndices(len(expected))

	for i, v := range got {
		if v != expected[i] {
			t.Errorf("Expected %d at index %d, got %d\n", expected[i], i, v)
		}
	}
}

func Test_BitReverseSampleOrder(t *testing.T) {
	input := []complex128{0, 1, 2, 3, 4, 5, 6, 7}
	expected := []complex128{0, 4, 2, 6, 1, 5, 3, 7}

	BitReverseSampleOrder(input)

	test.ValidateFFT(t, expected, input)
}

func Test_IterativeFFT(t *testing.T) {
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
			fft := make([]complex128, len(tc.samples))
			copy(fft, tc.samples)

			IterativeFFT(fft)

			test.ValidateFFT(t, tc.expected, fft)
		})
	}
}
