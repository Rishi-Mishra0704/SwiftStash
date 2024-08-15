clean:
	rm -rf ./bin/*

build: clean
	go build -o ./bin/SwiftStash main.go

run: build
	./bin/SwiftStash

PHONY: clean build run