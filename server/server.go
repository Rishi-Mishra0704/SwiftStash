package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/Rishi-Mishra0704/SwiftStash/cache"
	"github.com/Rishi-Mishra0704/SwiftStash/client"
	"github.com/Rishi-Mishra0704/SwiftStash/cmd"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
	LeaderAddr string
}

type Server struct {
	ServerOpts
	Cache   cache.Cacher
	members map[*client.Client]struct{}
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		Cache:      c,
		members:    make(map[*client.Client]struct{}),
	}

}

// Start the server
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}

	if !s.IsLeader && len(s.LeaderAddr) != 0 {
		go func() {
			if err := s.dialLeader(); err != nil {
				log.Println(err)
			}
		}()
	}

	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(conn)
	}
}

// handleConn handles incoming connections
func (s *Server) handleConn(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	//fmt.Println("connection made:", conn.RemoteAddr())

	for {
		cmd, err := cmd.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}
		go s.HandleCommand(conn, cmd)
	}

	// fmt.Println("connection closed:", conn.RemoteAddr())
}

func (s *Server) HandleCommand(conn net.Conn, command any) {
	switch v := command.(type) {
	case *cmd.CommandSet:
		_ = s.handleSetCommand(conn, v)
	case *cmd.CommandGet:
		_ = s.handleGetCommand(conn, v)
	case *cmd.CommandJoin:
		_ = s.handleJoinCommand(conn, v)
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
	// log.Printf("SET %s to %s", string(command.Key), string(command.Value))
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
func (s *Server) dialLeader() error {
	conn, err := net.Dial("tcp", s.LeaderAddr)
	if err != nil {
		return fmt.Errorf("failed to dial leader [%s]", s.LeaderAddr)
	}

	log.Println("connected to leader:", s.LeaderAddr)

	if err = binary.Write(conn, binary.LittleEndian, cmd.CmdJoin); err != nil {
		return err
	}

	s.handleConn(conn)

	return nil
}

func (s *Server) handleJoinCommand(conn net.Conn, _ *cmd.CommandJoin) error {
	fmt.Println("member just joined the cluster:", conn.RemoteAddr())

	s.members[client.NewFromConn(conn)] = struct{}{}

	return nil
}
