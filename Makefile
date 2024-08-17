clean:
	rm -rf ./bin/*

build: clean
	go build -o ./bin/SwiftStash main.go

run: build
	./bin/SwiftStash


runFollower: build
	./bin/SwiftStash --listenAddr :4000 --leaderAddr :3000 

test:
	@go test -v ./...


PHONY: clean build run runFollower test