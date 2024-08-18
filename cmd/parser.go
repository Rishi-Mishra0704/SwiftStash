package cmd

// CommandParser defines methods for parsing commands and converting them to bytes
type CommandParser interface {
	Bytes() []byte
}
