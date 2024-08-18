package client

import (
	"context"
	"fmt"
	"net"

	"github.com/Rishi-Mishra0704/SwiftStash/cmd"
)

// Options represents the client options
type Options struct{}

// Client represents a client connection
type Client struct {
	Conn net.Conn
}

// NewClient creates a new client connection
func NewClient(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)

	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn}, nil
}

// Set sends a SET command to the server
func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) error {
	command := &cmd.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.Conn.Write(command.Bytes())
	if err != nil {
		return err
	}
	resp, err := cmd.ParseSetResponse(c.Conn)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	if resp.Status != cmd.StatusOK {
		return fmt.Errorf("server returned error: %s", resp.Status.String())
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	command := &cmd.CommandGet{
		Key: key,
	}

	_, err := c.Conn.Write(command.Bytes())
	if err != nil {
		return nil, err
	}
	resp, err := cmd.ParseGetResponse(c.Conn)
	if err != nil {
		return nil, err
	}
	if resp.Status == cmd.StatusKeyNotFound {
		return nil, fmt.Errorf("could not find key (%s)", key)
	}
	if resp.Status != cmd.StatusOK {
		return nil, fmt.Errorf("server returned error: %s", resp.Status.String())
	}

	return resp.Value, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.Conn.Close()
}
