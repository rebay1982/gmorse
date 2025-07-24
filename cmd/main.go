package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
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

	// List playback devices (not interested in these.)
	fmt.Println("Playback devices:")
	pDevs, err := ctx.Devices(malgo.Playback)
	for i, dev := range pDevs {
		fmt.Printf("%d: %s\n", i, dev.Name())

	}

	// List capture devices
	fmt.Println("")
	fmt.Println("Capture devices:")
	cDevs, err := ctx.Devices(malgo.Capture)
	for i, cDev := range cDevs {
		fmt.Printf("%d: %s, default: %d\n", i, cDev.Name(), cDev.IsDefault)
	}

	fmt.Println("\n\n--- Initializing capture on default device ---")
	// Setup device to validate capture.
	const (
		sampleRate = 44100
		blockSize  = 2048
		toneFreq   = 700.0
	)

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	selInfo := cDevs[6]
	deviceId := selInfo.ID

	deviceConfig.Capture.DeviceID = deviceId.Pointer()
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1

	////buffer := make([]byte, blockSize << 1)
	// Avoid recreating these every time the onReceiveFrames function is called.
	samples := make([]float64, blockSize)
	freqSpectrum := make([]complex128, blockSize)
	magnitudes := make([]float64, blockSize)
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
			freqSpectrum[i] = complex(s, 0)
		}

		fft.IterativeFFT(freqSpectrum)

		hannFactorRMS := windowing.HannFactorRMS(int(sampleCount))
		halfSampleCount := int(sampleCount >> 1)
		for i := range halfSampleCount {
			// Normalize magnitudes and take into account the hanning window that was applied on the input samples before FFT.
			magnitudes[i] = fft.ComputeMagnitude(freqSpectrum[i]) / hannFactorRMS
		}

		timeDiff := time.Now().Sub(startTime)
		fmt.Printf("Processed %d in %d us          \n", sampleCount, timeDiff/time.Microsecond)

		for j := 9; j >= 0; j-- {
//			for i := range halfSampleCount {
//				meter := int(magnitudes[i] * 10)
			if int(magnitudes[1] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[4] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[11] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[29] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[74] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[146] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[232] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[329] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[557] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

			if int(magnitudes[788] * 10) > j {
				fmt.Print("::")
			} else {
				fmt.Print("  ")
			}

//				if meter >= j {
//					fmt.Print("::")
//				} else {
//					fmt.Print("  ")
//				}
//			}
			fmt.Println()
		}

		// Bring the cursor 10 lines up.
		fmt.Print("\033[11A\r")
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
}
