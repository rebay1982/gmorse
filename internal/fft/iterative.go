package fft

import (
	"math"
)

func IterativeFFT(samples []complex128) {
	BitReverseSampleOrder(samples)
	a := samples
	n := len(a)

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
				u := a[start+j]
				t := w * a[start+j+halfSize]

				// Butterfly outputs
				a[start+j] = u + t
				a[start+j+halfSize] = u - t

				// Update twiddle factor
				w *= wm
			}
		}
	}
}

func BitReverseSampleOrder(samples []complex128) {
	nbSamples := len(samples)
	indices := preComputeBitReverseIndices(nbSamples)

	// Half size only works because the sample size will always be a power of two.
	for i := range nbSamples >> 1 {
		j := indices[i]

		// Don't use arithmetic swaps, costly to compute and can introduce rounding errors.
		tmp := samples[i]
		samples[i] = samples[j]
		samples[j] = tmp
	}
}

// preComputeBitReverseIndices call once to get the bit reverse indices for n.
func preComputeBitReverseIndices(n int) []int {
	fn := float64(n)
	maskSize := uint(math.Log2(fn))

	indices := make([]int, n)
	for i := range uint(n) {
		indices[i] = int(reverseBits(i, maskSize))
	}

	return indices
}

func reverseBits(in, k uint) uint {
	var out uint

	for i := uint(0); i < k; i++ {
		out = (out << 1) | (in & 1)
		in = in >> 1
	}

	return out
}
