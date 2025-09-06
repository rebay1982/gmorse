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
// Setup device to validate capture.
const (
	sampleRate = 8000
	periodSizeMS = 10
	blockSize  = 128
	toneFreq   = 500.0
)

// Avoid recreating these every time the onReceiveFrames function is called.
var samples []float64 = make([]float64, blockSize)
var frequencies []float64 = []float64{500, 550, 600, 650, 700, 750, 800, 850, 900, 950}
var mags []float64 = make([]float64, len(frequencies))
var detectionCh = make(chan bool, 1)

//var onReceiveStart = time.Now()
func OnReceiveFrames(_, iSamples[]byte, sampleCount uint32) {
	//onReceiveEnd := time.Now()
	//fmt.Printf("Call took: %d\n", onReceiveEnd.Sub(onReceiveStart).Milliseconds())
	//onReceiveStart = time.Now()

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

	// Display detection
	detection := false
	for i := range frequencies {
		if mags[i] > 10.0 {
			detection = true
		}
	}

	detectionCh <- detection

}

var edgeDetected = false
var prevDetection = false

var onStart time.Time
var onEnd time.Time
var onLength time.Duration

var offStart time.Time
var offEnd time.Time
var offLength time.Duration

func Calibration(detectionLength <- chan int64, done <- chan struct{}) {


}



func HandleDetection(detectionCh <- chan bool, done <- chan struct{}) {
	calibrating := true
	calibrationSamplesCount := 0
	noEdgeCount := 0
	noEdgeCalibrationTrigger := 1500 // 1500 * 10 (10 ms sample window) == 15seconds
	maxCalibrationSamples := 10
	calibrationSamples := make([]time.Duration, 10)

	for {
		select {
		case morseDetected := <-detectionCh:
			// Edge Detection
			if morseDetected != prevDetection {
				noEdgeCount = 0		// Reset no edge count for calibration trigger.



				// Rising edge
				if morseDetected {
					offEnd = time.Now()
					onStart = time.Now()
				// Falling edge
				} else {
					onEnd = time.Now()
					offStart = time.Now()
				}



				if onEnd.After(onStart) {
					delta := onEnd.Sub(onStart)
					// We're still calibrating.
					if calibrating {
						if calibrationSamplesCount >= maxCalibrationSamples {
							// Determine WPM based on collected samples.

							// Exit Calibration
							calibrationSamplesCount = 0
							calibrating = false

						} else {
							calibrationSamples[calibrationSamplesCount] = delta
						}
					}
					fmt.Printf("Last on lasted %d\n", delta.Milliseconds())
				}



				prevDetection = morseDetected

			// No edge detection
			} else {
				noEdgeCount++
				if noEdgeCount > noEdgeCalibrationTrigger {
					noEdgeCount = 0
					calibrating = true
				}
			}


		case <- done:
			fmt.Println("Exiting detection routine")
			return 
		}
	}
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
	deviceConfig.PeriodSizeInMilliseconds = periodSizeMS
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


	done := make(chan struct{}, 1)
  go HandleDetection(detectionCh, done)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	close(done)

	fmt.Println("\nExiting...")
}

