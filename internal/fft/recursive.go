package fft

import (
	"fmt"
	"math"
)

func RecursiveFFT(samples []complex128) ([]complex128, error) {
	n := len(samples)
	nHalved := n >> 1

	if !isPowerOfTwo(n) {
		return samples, fmt.Errorf("Input sample size must be a power of two.")
	}

	// We're at the last stage. At this stage, there's no tiddle to compute in the discreet fourier transform since
	//   m = 0..0, e(-i*2*pi*k*m/n) will be 1.
	if n == 1 {
		return samples, nil
	}

	// Split samples into odd and even samples (purely based on their indices)
	evenSamples := make([]complex128, nHalved)
	oddSamples := make([]complex128, nHalved)
	for i := range nHalved {
		evenSamples[i] = samples[i<<1]
		oddSamples[i] = samples[(i<<1)+1]
	}

	// Recursively compute the DFT of each even and odd set.
	evenFFT, _ := RecursiveFFT(evenSamples)
	oddFFT, _ := RecursiveFFT(oddSamples)

	// Finally compute the DFT.
	out := make([]complex128, n)
	for k := range nHalved {

		// Twiddle factor's angle -- constant for a single frequency.
		angle := -2 * math.Pi * float64(k) / float64(n)
		twiddle := complex(math.Cos(angle), math.Sin(angle)) // Compute the twiddle factor.

		out[k] = evenFFT[k] + twiddle*oddFFT[k]
		out[k+nHalved] = evenFFT[k] - twiddle*oddFFT[k]
	}

	return out, nil
}
