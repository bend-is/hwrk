package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var ErrNotConnected = errors.New("not connected")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	conn net.Conn
	in   io.ReadCloser
	out  io.Writer

	addr    string
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		in:      in,
		out:     out,
		addr:    address,
		timeout: timeout,
	}
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *client) Close() error {
	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.Close()
}

func (c *client) Send() error {
	if c.conn == nil {
		return ErrNotConnected
	}

	scanner := bufio.NewScanner(c.in)
	for scanner.Scan() {
		b := append(scanner.Bytes(), '\n')
		if _, err := c.conn.Write(b); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (c *client) Receive() error {
	if c.conn == nil {
		return ErrNotConnected
	}

	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		b := append(scanner.Bytes(), '\n')
		if _, err := c.out.Write(b); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}

	return nil
}
