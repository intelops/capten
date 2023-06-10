.PHONY: build

build:
	go mod tidy
	go mod vendor
	go build -o capten ./cmd/main.go 
