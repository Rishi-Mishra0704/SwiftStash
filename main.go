package main

import (
	"log"
	"net"
	"time"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/server"
)

func main() {
	opts := server.ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	// For testing purposes only,  uncomment the following lines
	// to simulate a client connection
	go func() {
		time.Sleep(5 * time.Second)
		conn, err := net.Dial("tcp", ":3000")
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET Foo Bar 2500"))
	}()
	s := server.NewServer(opts, cache.NewCache())
	err := s.Start()
	if err != nil {
		return
	}
}
