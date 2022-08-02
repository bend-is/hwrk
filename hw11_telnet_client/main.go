package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const defaultTimeout = time.Second * 10

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", defaultTimeout, "connect timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		exitWithError("not enough arguments", true)
	}

	host := flag.Args()[0]
	port := flag.Args()[1]

	telnetCl := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	if err := telnetCl.Connect(); err != nil {
		exitWithError("connection failed: "+err.Error(), false)
	}

	go exitOnInterrupt()

	go func() {
		if err := telnetCl.Receive(); err != nil {
			exitWithError("receive message failed: "+err.Error(), false)
		}
	}()

	if err := telnetCl.Send(); err != nil {
		exitWithError("send message failed: "+err.Error(), false)
	}

	if err := telnetCl.Close(); err != nil {
		exitWithError("close connection failed: "+err.Error(), false)
	}

	fmt.Fprint(os.Stderr, "Exit")
}

func exitWithError(msg string, printUsage bool) {
	fmt.Println("error:", msg)

	if printUsage {
		fmt.Printf("\nusage: %s [-timeout] host port\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	os.Exit(1)
}

func exitOnInterrupt() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	<-exitChan

	os.Exit(0)
}
