package windowing

import "math"

// Hann Apply a Hanning window to a set of samples to limit and keep spectral leakage under control.
//
//	See https://en.wikipedia.org/wiki/Hann_function
func Hann(samples []float64) {
	N := len(samples)
	doublePie := 2 * math.Pi

	for i := range N {
		samples[i] *= 0.5 * (1 - math.Cos(doublePie*float64(i)/float64(N-1)))
	}
}
