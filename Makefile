.DEFAULT_GOAL := build

fmt:
	go fmt ./...

vet:
	go vet ./...

build: vet
	go build -o gmorse ./cmd/main.go

run: build
	./gmorse 2>/dev/null

test: build
	go test -v -count=1 ./...
