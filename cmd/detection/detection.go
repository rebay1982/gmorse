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
	"github.com/rebay1982/gdsp/filters"
	"github.com/rebay1982/gdsp/windowing"
	"github.com/rebay1982/gmorse/internal/decode"
)

// Setup device to validate capture.
const (
	sampleRate   = 8000
	periodSizeMS = 10
	blockSize    = 128
)

// Avoid recreating these every time the onReceiveFrames function is called.
var samples []float64 = make([]float64, blockSize)
var frequencies []float64 = []float64{500, 550, 600, 650, 700, 750, 800, 850, 900, 950}
var mags []float64 = make([]float64, len(frequencies))
var detectionCh = make(chan bool, 1)

func OnReceiveFrames(_, iSamples []byte, sampleCount uint32) {
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
	//maxMag := 0.0
	detection := false
	for i := range frequencies {
		//if mags[i] > maxMag {
		//	maxMag = mags[i]
		//}
		if mags[i] > 1.0 {
			detection = true
		}
	}

	//fmt.Println(maxMag)
	detectionCh <- detection
}

var prevDetection = false
var onStart time.Time
var onEnd time.Time
var onLength time.Duration
var offStart time.Time
var offEnd time.Time
var offLength time.Duration

func HandleDetection(detectionCh <-chan bool, decodeIn chan<- decode.Detection, done <-chan struct{}) {
	for {
		select {
		case morseDetected := <-detectionCh:

			// Edge Detection
			if morseDetected != prevDetection {
				// Rising edge
				if morseDetected {
					offEnd = time.Now()
					onStart = time.Now()

					offLength = offEnd.Sub(offStart)

					// Falling edge
				} else {
					onEnd = time.Now()
					offStart = time.Now()

					onLength = onEnd.Sub(onStart)
				}

				detection := decode.Detection{
					State: prevDetection,
				}
				if prevDetection {
					detection.Duration = onLength
				} else {
					detection.Duration = offLength
				}
				decodeIn <- detection

				prevDetection = morseDetected
			} else {
				// Time out after a second or two of silence.
				if !morseDetected {
					if time.Now().Sub(offStart).Seconds() > 2 {
						offEnd = time.Now()

						detection := decode.Detection{
							State:    prevDetection,
							Duration: offEnd.Sub(offStart),
						}
						offStart = time.Now()

						decodeIn <- detection
					}
				}
			}

		case <-done:
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

	fmt.Print("Configuring device and callback routine... ")
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

	// Initialize detection handling.
	done := make(chan struct{}, 1)
	decodeIn := make(chan decode.Detection)
	go HandleDetection(detectionCh, decodeIn, done)
	fmt.Println("Done")

	fmt.Print("Initializing morse decoder... ")
	decodeConfig := decode.DecoderConfig{
		Wpm:      25,
		Tolerace: 0.4,
	}
	decodeOut := make(chan string)
	decoder := decode.NewMorseDecoder(decodeIn, decodeOut, done, decodeConfig)
	decoder.StartDecode()
	fmt.Println("Done")

	fmt.Println("Starting capture...")
	device.Start()

	go func() {
		for {
			select {
			//case <-decodeOut:
			case msg := <-decodeOut:
				fmt.Print(msg)
			case <-done:
				return
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	close(done)

	fmt.Println("\nExiting...")
}
