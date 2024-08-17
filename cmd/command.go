package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CommandParser defines the interface for parsing and serializing cache commands.
type CommandParser interface {
	ParseMessage([]byte) (*Message, error)
	ToBytes() []byte
}

// Command represents cache commands like "SET" and "GET".
type Command string

const (
	CMDSET Command = "SET"
	CMDGET Command = "GET"
)

// MessageSet holds the data for a "SET" command.
type MessageSet struct {
	Key   []byte
	Value []byte
	TTL   time.Duration
}

// MessageGet holds the data for a "GET" command.
type MessageGet struct {
	Key []byte
}

// Message represents a parsed command message.
type Message struct {
	Command Command
	Key     []byte
	Value   []byte
	TTL     time.Duration
}

// ParseMessage parses a raw command into a Message struct.
func ParseMessage(rawCmd []byte) (*Message, error) {
	rawCmdStr := string(rawCmd)
	parts := strings.Split(rawCmdStr, " ")

	if len(parts) == 0 {
		return nil, errors.New("invalid protocol format")
	}

	msg := &Message{
		Command: Command(parts[0]),
		Key:     []byte(parts[1]),
	}

	if msg.Command == CMDSET {
		if len(parts) != 4 {
			return nil, errors.New("invalid SET command")
		}
		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid SET TTL")
		}
		msg.Value = []byte(parts[2])
		msg.TTL = time.Duration(ttl)
	}

	return msg, nil
}

// ToBytes converts a Message struct back into a byte slice.
func (m *Message) ToBytes() []byte {
	switch m.Command {
	case CMDSET:
		cmd := fmt.Sprintf("%s %s %s %d", m.Command, m.Key, m.Value, m.TTL)
		return []byte(cmd)
	case CMDGET:
		cmd := fmt.Sprintf("%s %s", m.Command, m.Key)
		return []byte(cmd)
	default:
		panic("invalid command")
	}
}
