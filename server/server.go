package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/cmd"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	Cache     cache.Cacher
	Followers map[net.Conn]struct{}
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		Cache:      c,
		// TODO: Only allocate followers if this is a leader
		Followers: make(map[net.Conn]struct{}),
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

	if !s.IsLeader {
		go func() {
			// Connect to leader
			conn, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Printf("Connected to leader at %s\n", s.LeaderAddr)
			if err != nil {
				log.Fatal(err)
			}
			s.handleConn(conn)
		}()

	}

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
	fmt.Println("Connected to ", conn.RemoteAddr().String())

	if s.IsLeader {
		s.Followers[conn] = struct{}{}
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error reading from connection: %s", err)
			break
		}
		msg := buf[:n]
		go s.handleCMD(conn, msg)
	}

}

func (s *Server) handleCMD(conn net.Conn, rawCmd []byte) {

	var (
		msg *cmd.Message
		err error
	)

	msg, err = cmd.ParseMessage(rawCmd)
	if err != nil {
		log.Printf("Error parsing command: %s", err)
		return
	}

	fmt.Printf("recieved command %s\n", msg.Command)

	switch msg.Command {
	case cmd.CMDSET:
		err = s.handleSET(conn, msg)
	case cmd.CMDGET:
		err = s.handleGET(conn, msg)
	}

	if err != nil {
		log.Printf("Error handling command: %s", err)
		conn.Write([]byte("ERR: " + err.Error()))
	}

}

func (s *Server) handleSET(conn net.Conn, msg *cmd.Message) error {

	err := s.Cache.Set(msg.Key, msg.Value, msg.TTL)
	if err != nil {
		return err
	}
	go s.sendToFollower(context.TODO(), msg)

	return nil
}

func (s *Server) handleGET(conn net.Conn, msg *cmd.Message) error {
	val, err := s.Cache.Get(msg.Key)
	if err != nil {
		return err
	}
	_, err = conn.Write(val)

	return err
}
func (s *Server) sendToFollower(ctx context.Context, msg *cmd.Message) error {
	for conn := range s.Followers {
		fmt.Println("Sending key to followers")
		rawMsg := msg.ToBytes()
		fmt.Println("Sending raw msg to followers", string(rawMsg))
		_, err := conn.Write(rawMsg)
		if err != nil {
			fmt.Printf("Error writing to follower: %s", err)
			continue
		}
	}
	return nil
}
