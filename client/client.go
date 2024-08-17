package client

import (
	"context"
	"net"

	"github.com/Rishi-Mishra0704/SwiftStash/cmd"
)

type Options struct{}

type Client struct {
	Conn net.Conn
}

func NewClient(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)

	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn}, nil
}

func (c *Client) Set(ctx context.Context, key, value []byte, ttl int) (any, error) {
	command := &cmd.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	_, err := c.Conn.Write(command.Bytes())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}
