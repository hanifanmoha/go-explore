.PHONY: build run run-dev clean

build:
	go build -o bin/main main.go

run: build
	./bin/main

run-dev:
	air

clean:
	rm -rf bin