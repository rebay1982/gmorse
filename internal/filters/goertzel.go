package goertzel

import "math"

func Goertzel(sampleRate, targetFreq float64, samples []float64) complex128 {
	N := float64(len(samples))
	k := int(math.Round(N * targetFreq / sampleRate))
	w := 2 * math.Pi * float64(k) / N
	cosine := math.Cos(w)
	sine := math.Sin(w)

	var sPrev, sPrevPrev float64

	// Step 1: Accumulate "energy" into target frequency.
	for _, x := range samples {
		s := x + 2*cosine*sPrev - sPrevPrev
		sPrevPrev = sPrev
		sPrev = s
	}

	// Step 2: Compute real and imaginary components.
	real := sPrev - cosine*sPrevPrev
	img := sine * sPrevPrev

	return complex(real, img)
}
