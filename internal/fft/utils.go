package fft

import (
	"math"
)

// Industry standard is to divide by 32768.
//
//	The max positive 16-bit PCM value is 32767, normalized will be 0.99997.
//	The min negative 16-bit PCM value is -32768, normalized will be -1.0.
//
//	If we divide by 32767 so that the max positive is 1.0, the min negative will be 1.00003 which throws off
//	  normalization.
const PCM_16_DIVISOR = 32768.0

func NormalizePCM16Samples(samples []int16) []float64 {
	N := len(samples)
	normalized := make([]float64, N)

	for i, pcm := range samples {
		normalized[i] = float64(pcm) / PCM_16_DIVISOR

	}

	return normalized
}

func NormalizePCM16(sample int16) float64 {
	return float64(sample) / PCM_16_DIVISOR
}

func ComputeMagnitude(sample complex128) float64 {
	r := real(sample)
	i := imag(sample)

	return math.Sqrt((r * r) + (i * i))
}

func ComputePhaseRad(sample complex128) float64 {
	r := real(sample)
	i := imag(sample)

	return math.Atan(i / r)
}

func ComputePhaseDeg(sample complex128) float64 {
	return ComputePhaseRad(sample) * 180 / math.Pi
}

func isPowerOfTwo(x int) bool {
	if x == 0 {
		return false
	} else {
		return x&(x-1) == 0
	}
}
