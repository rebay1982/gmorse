.DEFAULT_GOAL := build

fmt:
	go fmt ./...

vet:
	go vet ./...

build-spectrum: vet
	go build -o spectrum ./cmd/spectrum/spectrum.go

build-goertzel: vet
	go build -o goertzel ./cmd/goertzel/goertzel.go

build-detection: vet
	go build -o detection ./cmd/detection/detection.go

spectrum: build-spectrum
	./spectrum 2>/dev/null

goertzel: build-goertzel
	./goertzel 2>/dev/null

detection: build-detection
	./detection 2>/dev/null

test: 
	go test -v -count=1 ./...
