package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CommandInvalid struct{}

func TestParseSetCommand(t *testing.T) {
	command := &CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTL:   2,
	}
	fmt.Println(command.Bytes())
	r := bytes.NewReader(command.Bytes())
	pCmd, err := ParseCommand(r)
	assert.Nil(t, err)
	assert.Equal(t, command, pCmd)
}
func TestParseGetCommand(t *testing.T) {
	command := &CommandGet{
		Key: []byte("foo"),
	}
	fmt.Println(command.Bytes())
	r := bytes.NewReader(command.Bytes())
	pCmd, err := ParseCommand(r)
	assert.Nil(t, err)
	assert.Equal(t, command, pCmd)
}

func TestParseInvalidCommand(t *testing.T) {
	// Simulate an invalid command byte
	command := []byte{255}
	fmt.Println(command)
	r := bytes.NewReader(command)
	_, err := ParseCommand(r)
	assert.Error(t, err, "invalid command")
}

func BenchmarkParseCommand(b *testing.B) {
	command := &CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTL:   2,
	}
	r := bytes.NewReader(command.Bytes())
	for i := 0; i < b.N; i++ {
		ParseCommand(r)
	}
}
