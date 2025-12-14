package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	ErrFailedConnection = errors.New("unable to establish a connection")
	ErrConnectWasClosed = errors.New("connection was closed by the remote host")
	ErrContextTimeout   = errors.New("context deadline exceeded")
	ErrEOF              = errors.New("unexpected end of input")
	ErrOther            = errors.New("an unspecified error occurred")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type client struct {
	address    string
	conn       net.Conn
	in         *bufio.Scanner
	outScanner *bufio.Scanner
	out        io.Writer
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &client{
		ctx:     ctx,
		cancel:  cancel,
		address: address,
		in:      bufio.NewScanner(in),
		out:     out,
	}
}

func (t *client) Connect() (err error) {
	dialer := &net.Dialer{}
	if t.conn, err = dialer.DialContext(t.ctx, "tcp", t.address); err != nil {
		t.cancel()
		return ErrFailedConnection
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.address)
	t.outScanner = bufio.NewScanner(t.conn)
	return nil
}

func (t *client) Close() (err error) {
	defer t.cancel()
	if err = t.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (t *client) Send() (err error) {
	select {
	case <-t.ctx.Done():
		return ErrContextTimeout
	default:
		if !t.in.Scan() {
			if t.in.Err() == nil {
				fmt.Fprint(os.Stderr, "....EOF\n")
				t.Close()
				return ErrEOF
			}
			return ErrOther
		}
		if _, err := t.conn.Write(fmt.Appendf(nil, "%s\n", t.in.Text())); err != nil {
			return ErrConnectWasClosed
		}
	}
	return nil
}

func (t *client) Receive() (err error) {
	select {
	case <-t.ctx.Done():
		return ErrContextTimeout
	default:
		if !t.outScanner.Scan() {
			fmt.Fprint(os.Stderr, "...Connection was closed by peer\n")
			t.cancel()
			return ErrConnectWasClosed
		}
		fmt.Fprintf(t.out, "%s\n", t.outScanner.Text())
	}
	return nil
}
