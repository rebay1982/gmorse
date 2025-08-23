package main

import (
	"encoding/binary"
	"fmt"
	//"math"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/rebay1982/gdsp/fft"
	"github.com/rebay1982/gdsp/windowing"
	"github.com/rebay1982/gdsp/filters"
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

	selectedDevId, err := strconv.Atoi(strDevId)
	if err != nil {
		fmt.Println("Bad user input, exiting...")
		os.Exit(1)
	}
	if selectedDevId < 0 || selectedDevId >= len(cDevs) {
		fmt.Printf("Bad user input, expecting value between 0 and %d, got %d\n", len(cDevs), selectedDevId)
		os.Exit(1)
	}

	fmt.Println("\n\n--- Initializing capture on default device ---")
	// Setup device to validate capture.
	const (
		sampleRate = 8000
		blockSize  = 256
		toneFreq   = 500.0
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
	frequencies := []float64{500, 550, 600, 650, 700, 750, 800, 850, 900, 950}
	mags := make([]float64, len(frequencies))
	onReceiveFrames := func(_, iSamples []byte, sampleCount uint32) {
		startTime := time.Now()

		// Normalize
		for i := range sampleCount {
			samples[i] = fft.NormalizePCM16(int16(binary.LittleEndian.Uint16(iSamples[i<<1 : (i+1)<<1])))
		}

		// Window (reduces spectral leakage)
		// Only apply it to the samples, not the padding.
		windowing.Hann(samples[:sampleCount])

		// Retrieve Goertzel calculation for all frequencies.
		for i, f := range frequencies {
			goertzel := filters.Goertzel(sampleRate, f, samples)
			mags[i] = fft.ComputeMagnitude(goertzel) * 2 // Compensate for the Hanning window
		}

		timeDiff := time.Now().Sub(startTime)
		fmt.Printf("Processed %d in %d us          \n", sampleCount, timeDiff/time.Microsecond)

		// Display detection
		for i, f := range frequencies {
			if mags[i] > 15.0 {
				fmt.Printf("%.f: DETECTION -- %02.2f         \n", f, mags[i])
			} else {
				fmt.Printf("%.f:           -- %02.2f         \n", f, mags[i])
			}
		}

		fmt.Printf("\033[%dA\r", len(frequencies) + 1)
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
