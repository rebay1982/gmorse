package fft

import (
	"testing"
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

	validateIterativeFFT(t, expected, input)
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fft := make([]complex128, len(tc.samples))
			copy(fft, tc.samples)

			IterativeFFT(fft)

			validateIterativeFFT(t, tc.expected, fft)
		})
	}
}

func validateIterativeFFT(t *testing.T, expected, fft []complex128) {
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
