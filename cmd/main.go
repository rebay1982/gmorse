package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/rebay1982/gmorse/internal/fft"
	"github.com/rebay1982/gmorse/internal/windowing"
)

func main() {

	// Setup malgo
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		fmt.Println("Did not work...")
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	// List capture devices
	fmt.Println("")
	fmt.Println("Capture devices:")
	cDevs, err := ctx.Devices(malgo.Capture)
	for i, cDev := range cDevs {
		fmt.Printf("%d: %s, default: %d\n", i, cDev.Name(), cDev.IsDefault)
	}

	fmt.Println("-- Select input device: ")
	var strDevId string
	fmt.Scanln(&strDevId)

	i, err := strconv.Atoi(strDevId)
	if err != nil {
		fmt.Println("Bad user input, exiting...")
		os.Exit(1)
	}
	if i < 0 || i >= len(cDevs) {
		fmt.Printf("Bad user input, expecting value between 0 and %d, got %d\n", len(cDevs), i)
		os.Exit(1)
	}

	fmt.Println("\n\n--- Initializing capture on default device ---")
	// Setup device to validate capture.
	const (
		sampleRate = 8000
		blockSize  = 256
		toneFreq   = 700.0
	)

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	selInfo := cDevs[1]
	deviceId := selInfo.ID

	deviceConfig.Capture.DeviceID = deviceId.Pointer()
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1

	// Avoid recreating these every time the onReceiveFrames function is called.
	samples := make([]float64, blockSize)
	fspec := make([]complex128, blockSize)
	normalizedMags := make([]float64, blockSize)
	onReceiveFrames := func(_, iSamples []byte, sampleCount uint32) {
		startTime := time.Now()

		// Normalize
		for i := range sampleCount {
			samples[i] = fft.NormalizePCM16(int16(binary.LittleEndian.Uint16(iSamples[i<<1 : (i+1)<<1])))
		}

		// Window (reduces spectral leakage)
		// Only apply it to the samples, not the padding.
		windowing.Hann(samples[:sampleCount])

		for i, s := range samples {
			fspec[i] = complex(s, 0)
		}

		//freqSpectrum := fft.RecursiveFFT(fspec)
		fft.IterativeFFT(fspec)
		freqSpectrum := fspec

		hannFactorRMS := windowing.HannFactorRMS(int(sampleCount))
		halfSampleCount := int(sampleCount >> 1)
		for i := 1; i < int(sampleCount)-1; i++ {
			// Normalize magnitudes and take into account the hanning window that was applied on the input samples before FFT.
			normalizedMags[i] = 2.0 * (fft.ComputeMagnitude(freqSpectrum[i]) / float64(sampleCount)) / hannFactorRMS
		}
		normalizedMags[0] = (fft.ComputeMagnitude(freqSpectrum[0]) / float64(sampleCount)) / hannFactorRMS
		normalizedMags[halfSampleCount-1] = (fft.ComputeMagnitude(freqSpectrum[halfSampleCount-1]) / float64(sampleCount)) / hannFactorRMS

		timeDiff := time.Now().Sub(startTime)
		fmt.Printf("Processed %d in %d us          \n", sampleCount, timeDiff/time.Microsecond)

		// Display spectrum
		for j := 0.0; j < 10.0; j += 0.5 {
			dbFloor := j * -10.0
			fmt.Printf("%06.2f ", dbFloor)

			for i := range 100 {
				mag := 20 * math.Log10(normalizedMags[i])

				if mag > dbFloor {
					fmt.Print("::")
				} else {
					fmt.Print("  ")
				}
			}
			fmt.Println()
		}

		// Bring the cursor 10 lines up.
		fmt.Print("\033[21A\r")
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Data: onReceiveFrames,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	defer device.Uninit()

	fmt.Println("\n\n--- Initializing capture on default device ---")
	device.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("\nExiting...")
}
