package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/client"
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
		client, err := client.NewClient(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < 10; i++ {
			SendCommand(client)
			time.Sleep(200 * time.Millisecond)

		}
		client.Close()
		time.Sleep(1 * time.Second)
	}()
	s := server.NewServer(opts, cache.NewCache())
	err := s.Start()
	if err != nil {
		return
	}
}

func SendCommand(c *client.Client) {
	command := &cmd.CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTL:   0,
	}

	_, err := c.Set(context.Background(), command.Key, command.Value, command.TTL)
	if err != nil {
		log.Fatal(err)
	}

}
