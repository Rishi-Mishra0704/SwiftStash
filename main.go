package main

import (
	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/server"
)

func main() {
	opts := server.ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	s := server.NewServer(opts, cache.NewCache())
	err := s.Start()
	if err != nil {
		return
	}
}
