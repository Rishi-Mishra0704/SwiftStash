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
	CmdNone Command = iota
	CmdSET
	CmdGET
	CmdDel
)

// Status represents different statuses
type Status byte

const (
	StatusNone Status = iota
	StatusOK
	StatusError
	StatusKeyNotFound
)

// CommandSet represents the SET command
type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
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

// CommandGet represents the GET command
type CommandGet struct {
	Key []byte
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

// ParseGetCommand parses a GET command from the reader
func ParseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, cmd.Key)

	return cmd
}

// ResponseGet represents a response for the GET command
type ResponseGet struct {
	Status Status
	Value  []byte
}

// Bytes returns the byte representation of ResponseGet
func (r *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	valueLen := int32(len(r.Value))
	binary.Write(buf, binary.LittleEndian, valueLen)
	binary.Write(buf, binary.LittleEndian, r.Value)
	return buf.Bytes()
}

// ResponseSet represents a response for the SET command
type ResponseSet struct {
	Status Status
}

// Bytes returns the byte representation of ResponseSet
func (r *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, r.Status)
	return buf.Bytes()
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

func ParseSetResponse(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}

	err := binary.Read(r, binary.LittleEndian, &resp.Status)

	return resp, err
}

func ParseGetResponse(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}
	if err := binary.Read(r, binary.LittleEndian, &resp.Status); err != nil {
		return resp, err
	}

	var valueLen int32
	if err := binary.Read(r, binary.LittleEndian, &valueLen); err != nil {
		return resp, err
	}

	resp.Value = make([]byte, valueLen)
	if err := binary.Read(r, binary.LittleEndian, &resp.Value); err != nil {
		return resp, err
	}

	return resp, nil
}

func (s Status) String() string {

	switch s {
	case StatusOK:
		return "OK"
	case StatusError:
		return "Error"
	case StatusNone:
		return "None"

	case StatusKeyNotFound:
		return "Key Not Found"
	default:
		return "INVALID"
	}
}
