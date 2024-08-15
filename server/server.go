package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	Cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		Cache:      c,
	}

}

// Start the server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		log.Printf("Error starting server: %s", err)
	}
	defer ln.Close()

	log.Printf("Server listening on [%s]", s.ListenAddr)

	// handle connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		go s.handleConn(conn)

	}
}

// handleConn handles incoming connections
func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from connection: %s", err)
			break
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}

}
