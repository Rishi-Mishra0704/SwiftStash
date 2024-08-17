package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/cmd"
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

	go func() {
		time.Sleep(2 * time.Second)
		for i := 0; i < 10; i++ {
			SendCommand()
			time.Sleep(200 * time.Millisecond)

		}
	}()
	s := server.NewServer(opts, cache.NewCache())
	err := s.Start()
	if err != nil {
		return
	}
}

func SendCommand() {
	command := &cmd.CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTL:   0,
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(command.Bytes())
}
