package cmd

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSET Command = "SET"
	CMDGET Command = "GET"
)

type MessageSet struct {
	Key   []byte
	Value []byte
	TTL   time.Duration
}

type MessageGet struct {
	Key []byte
}

type Message struct {
	Command Command
	Key     []byte
	Value   []byte
	TTL     time.Duration
}

func ParseMessage(rawCmd []byte) (*Message, error) {

	var (
		rawCmdStr = string(rawCmd)
		parts     = strings.Split(rawCmdStr, " ")
	)
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
