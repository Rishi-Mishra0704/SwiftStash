package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/client"
	"github.com/Rishi-Mishra0704/SwiftStash/server"
)

func main() {
	var (
		listenAddr = flag.String("listenaddr", ":3000", "listen address of the server")
		leaderAddr = flag.String("leaderaddr", "", "listen address of the leader")
	)
	flag.Parse()

	opts := server.ServerOpts{
		ListenAddr: *listenAddr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 10)
		if opts.IsLeader {
			SendStuff()
		}
	}()

	server := server.NewServer(opts, cache.NewCache())
	_ = server.Start()
}

func SendStuff() {
	for i := 0; i < 100; i++ {
		go func(i int) {
			c, err := client.NewClient(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("val_%d", i))
			)

			err = c.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}

			fetchedValue, err := c.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(fetchedValue))

			_ = c.Close()
		}(i)
	}
}
