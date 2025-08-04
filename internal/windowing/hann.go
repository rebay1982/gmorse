package windowing

import "math"

const (
	DOUBLE_PIE = 2 * math.Pi
)

// Hann Apply a Hanning window to a set of samples to limit and keep spectral leakage under control.
//
//	See https://en.wikipedia.org/wiki/Hann_function
func Hann(samples []float64) {
	N := len(samples)

	for i := range N {
		samples[i] *= 0.5 * (1 - math.Cos(DOUBLE_PIE*float64(i)/float64(N-1)))
	}
}

func HannFactorRMS(sampleSize int) float64 {
	N := sampleSize
	var factorRMS float64 = 0.0

	for i := range N {
		w := 0.5 * (1 - math.Cos(DOUBLE_PIE*float64(i)/float64(N-1)))
		factorRMS += w * w
	}

	return math.Sqrt(factorRMS / float64(N))
}
