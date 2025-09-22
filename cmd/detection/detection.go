package main

import (
	"encoding/binary"
	"fmt"
	//"math"
	"os"
	"os/signal"
	"strconv"
	//"time"

	"github.com/gen2brain/malgo"
	"github.com/rebay1982/gdsp/fft"
	"github.com/rebay1982/gdsp/filters"
	"github.com/rebay1982/gdsp/windowing"
)

// Setup device to validate capture.
const (
	sampleRate = 8000
	blockSize  = 256
	toneFreq   = 500.0
)

// Avoid recreating these every time the onReceiveFrames function is called.
var samples []float64 = make([]float64, blockSize)
var frequencies []float64 = []float64{500, 550, 600, 650, 700, 750, 800, 850, 900, 950}
var mags []float64 = make([]float64, len(frequencies))

func OnReceiveFrames(_, iSamples []byte, sampleCount uint32) {
	//startTime := time.Now()

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

	//timeDiff := time.Now().Sub(startTime)
	//fmt.Printf("Processed %d in %d us          \n", sampleCount, timeDiff/time.Microsecond)

	// Display detection
	detection := false
	for i, _ := range frequencies {
		if mags[i] > 20.0 {
			detection = true
			//fmt.Printf("%.fHz: @@@@ -- %02.2f  \n", f, mags[i])
		}
		//else {
		//fmt.Printf("%.fHz:      -- %02.2f  \n", f, mags[i])
		//}
	}

	if detection {
		fmt.Print(".")
	} else {
		fmt.Print(" ")
	}

	//fmt.Printf("\033[%dA\r", len(frequencies) + 1)
}

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

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	selInfo := cDevs[1]
	deviceId := selInfo.ID

	deviceConfig.Capture.DeviceID = deviceId.Pointer()
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1

	captureCallbacks := malgo.DeviceCallbacks{
		Data: OnReceiveFrames,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	defer device.Uninit()

	fmt.Println("\n\n--- Initializing capture ---")
	device.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	fmt.Println("\nExiting...")
}
