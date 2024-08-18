package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

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
	defer conn.Close()

	fmt.Println("Connected to ", conn.RemoteAddr().String())

	for {
		cmd, err := cmd.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error parsing command: %s", err)
			break
		}
		go s.HandleCommand(conn, cmd)
	}

	fmt.Println("Disconnected from ", conn.RemoteAddr().String())
}

func (s *Server) HandleCommand(conn net.Conn, command any) {
	switch v := command.(type) {
	case *cmd.CommandGet:
		s.handleGetCommand(conn, v)
	case *cmd.CommandSet:
		s.handleSetCommand(conn, v)
	}
}

func (s *Server) handleGetCommand(conn net.Conn, command *cmd.CommandGet) error {

	resp := cmd.ResponseGet{}
	value, err := s.Cache.Get(command.Key)
	if err != nil {
		resp.Status = cmd.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = cmd.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return err
}

func (s *Server) handleSetCommand(conn net.Conn, command *cmd.CommandSet) error {
	log.Printf("SET %s to %s", string(command.Key), string(command.Value))
	resp := &cmd.ResponseSet{}
	if err := s.Cache.Set(command.Key, command.Value, time.Duration(command.TTL)); err != nil {
		resp.Status = cmd.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = cmd.StatusOK
	_, err := conn.Write(resp.Bytes())
	return err
}
