package fft

import (
	"fmt"
	"math"
)

// IterativeFFT Computes the FFT of samples, a complex128 slice, using the radix2 iterative method. This requires that
// samples matches a length of a power of 2.
func IterativeFFT(samples []complex128) error {
	if !isPowerOfTwo(len(samples)) {
		return fmt.Errorf("Input sample size must be a power of two.")
	}

	BitReverseSampleOrder(samples)
	n := len(samples)

	// Iterative FFT computation using butterfly operations
	for size := 2; size <= n; size <<= 1 {
		halfSize := size >> 1

		// Precompute twiddle factor increment for each step
		theta := -2 * math.Pi / float64(size)
		wm := complex(math.Cos(theta), math.Sin(theta))

		for start := 0; start < n; start += size {
			w := complex(1, 0)
			for j := range halfSize {
				// Butterfly inputs
				u := samples[start+j]
				t := w * samples[start+j+halfSize]

				// Butterfly outputs
				samples[start+j] = u + t
				samples[start+j+halfSize] = u - t

				// Update twiddle factor
				w *= wm
			}
		}
	}

	return nil
}

func BitReverseSampleOrder(samples []complex128) {
	nbSamples := len(samples)
	indices := preComputeBitReverseIndices(nbSamples)

	// Half size only works because the sample size will always be a power of two.
	for i := range nbSamples {
		j := indices[i]

		if i < j {
			// Don't use arithmetic swaps, costly to compute and can introduce rounding errors.
			samples[i], samples[j] = samples[j], samples[i]
		}
	}
}

// preComputeBitReverseIndices call once to get the bit reverse indices for n.
func preComputeBitReverseIndices(n int) []int {
	maskSize := math.Log2(float64(n))

	indices := make([]int, n)
	for i := range n {
		indices[i] = reverseBits(i, int(maskSize))
	}

	return indices
}

func reverseBits(in, k int) int {
	var out int

	for range k {
		out = (out << 1) | (in & 1)
		in = in >> 1
	}

	return out
}
