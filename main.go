package main

import (
	"flag"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/server"
)

func main() {

	listenAddr := flag.String("listenAddr", ":3000", "Listen address of the server")
	leaderAddr := flag.String("leaderAddr", "", "Listen address of the leader")
	flag.Parse()
	opts := server.ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	s := server.NewServer(opts, cache.NewCache())
	err := s.Start()
	if err != nil {
		return
	}
}
