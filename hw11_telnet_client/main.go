package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go-telnet [--timeout=10s] <host> <port>")
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := telnetClient.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "connection error:", err)
		return
	}
	defer telnetClient.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	sendDone := make(chan struct{})
	receiveDone := make(chan struct{})

	go func() {
		if err := telnetClient.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "send error:", err)
		}
		close(sendDone)
	}()

	go func() {
		if err := telnetClient.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "receive error:", err)
		}
		close(receiveDone)
	}()

	select {
	case <-sigCh:
		log.Println("signal received, closing connection")
	case <-sendDone:
	case <-receiveDone:
	}

	telnetClient.Close()
}
