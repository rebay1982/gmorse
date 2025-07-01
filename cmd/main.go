package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"

	"github.com/gen2brain/malgo"
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
		sampleRate = 8000
		blockSize  = 256
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
	onReceiveFrames := func(_, iSamples []byte, sampleCount uint32) {
		var maxAmplitude int16 = 0

		for i := range sampleCount {
			//for i := 0; i < int(sampleCount); i++ {
			amplitude := int16(binary.LittleEndian.Uint16(iSamples[i<<1 : (i+1)<<1]))
			if amplitude > maxAmplitude {
				maxAmplitude = amplitude
			}
		}

		fmt.Printf("Signal amplitude: %d             \r", maxAmplitude)
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
