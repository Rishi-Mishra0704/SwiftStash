package cmd

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Command byte

const (
	CmdNonce Command = iota
	CmdSET
	CmdGET
	CmdDel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int
}
type CommandGet struct {
	Key []byte
}

func ParseCommand(r io.Reader) any {
	var cmd Command

	binary.Read(r, binary.LittleEndian, &cmd)

	switch cmd {
	case CmdSET:
		return ParseSetCommand(r)
	case CmdGET:
		return ParseGetCommand(r)
	default:
		panic("Invalid Command")
	}
}

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
func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGET)
	keyLen := int32(len(c.Key))
	binary.Write(buf, binary.LittleEndian, keyLen)
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()

}

func ParseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	// Read the length of the key
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, cmd.Key)

	// Read the length of the value
	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, cmd.Value)

	// Read the TTL
	var ttl int64
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTL = int(ttl)

	return cmd
}
func ParseGetCommand(r io.Reader) *CommandGet {
	cmd := &CommandGet{}

	// Read the length of the key
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, cmd.Key)

	return cmd
}
