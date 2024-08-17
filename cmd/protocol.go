package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Command represents different types of commands
type Command byte

const (
	CmdNonce Command = iota
	CmdSET
	CmdGET
	CmdDel
)

// CommandParser defines methods for parsing commands and converting them to bytes
type CommandParser interface {
	Bytes() []byte
}

// CommandSet represents the SET command
type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}

// CommandGet represents the GET command
type CommandGet struct {
	Key []byte
}

// ParseCommand parses a command from the reader
func ParseCommand(r io.Reader) (CommandParser, error) {
	var cmd Command

	if err := binary.Read(r, binary.LittleEndian, &cmd); err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSET:
		return ParseSetCommand(r), nil
	case CmdGET:
		return ParseGetCommand(r), nil
	default:
		return nil, fmt.Errorf("invalid command")
	}
}

// Bytes returns the byte representation of CommandSet
func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSET)
	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)
	valueLen := int32(len(c.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, c.Value)
	binary.Write(buf, binary.LittleEndian, int64(c.TTL))
	return buf.Bytes()
}

// Bytes returns the byte representation of CommandGet
func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGET)
	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}

// ParseSetCommand parses a SET command from the reader
func ParseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, cmd.Value)

	var ttl int64
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)

	return cmd
}

// ParseGetCommand parses a GET command from the reader
func ParseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, cmd.Key)

	return cmd
}
