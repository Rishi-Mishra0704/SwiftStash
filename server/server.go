package server

import (
	"context"
	"log"
	"net"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/cmd"
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

func (s *Server) sendToFollower(ctx context.Context, msg *cmd.Message) error {
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
