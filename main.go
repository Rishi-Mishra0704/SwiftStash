package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/client"
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
		client, err := client.NewClient(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}
		err = client.Set(context.Background(), []byte("foo"), []byte("bar"), 0)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
		value, err := client.Get(context.Background(), []byte("foo"))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(200 * time.Millisecond)
		log.Printf("Value: %s", value)
		client.Close()
	}()
	s := server.NewServer(opts, cache.NewCache())
	err := s.Start()
	if err != nil {
		return
	}
}
