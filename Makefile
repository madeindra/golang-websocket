.PHONY:

setup:
	go mod tidy

build: .PHONY
	go build -o main ./cmd/main

run: .PHONY build
	./main