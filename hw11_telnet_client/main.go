package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "Connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("usage: go-telnet [--timeout=10s] <host> <port>")
	}
	address := net.JoinHostPort(args[0], args[1])

	telnetClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := telnetClient.Connect(); errors.Is(err, ErrFailedConnection) {
		log.Fatal(err.Error())
	}
	defer telnetClient.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("signal received, closing connection")
		telnetClient.Close()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			if err := telnetClient.Receive(); err != nil {
				break
			}
		}
		wg.Done()
	}()

	go func() {
		for {
			if err := telnetClient.Send(); err != nil {
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
