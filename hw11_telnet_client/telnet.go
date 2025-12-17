package main

import (
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *client) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	log.Printf("Connected to %s\n", t.address)
	return nil
}

func (t *client) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *client) Send() error {
	if _, err := io.Copy(t.conn, t.in); err != nil {
		return err
	}
	log.Println("EOF")
	return nil
}

func (t *client) Receive() error {
	if _, err := io.Copy(t.out, t.conn); err != nil {
		return err
	}
	log.Println("Connection was closed by peer")
	return nil
}
