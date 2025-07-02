package main

import (
	"flag"
	"fmt"
	"log"
	"tcp-scaner/internal/tcp"
	"time"
)

func StartScanner(data data) {
	scanner := tcp.Scanner{
		Address:     data.host,
		StartPort:   data.startPort,
		EndPort:     data.endPort,
		Concurrency: data.concurrency,
		Timeout:     time.Duration(data.timeout) * time.Millisecond,
	}

	if data.startPort < 1 || data.endPort > 65535 || data.startPort > data.endPort {
		fmt.Println("Invalid port range")
		return
	}

	log.Printf("Scanning %s ports %d–%d with %d workers and timeout %dms...\n",
		data.host, data.startPort, data.endPort, data.concurrency, data.timeout)

	openPorts := scanner.Run()

	if len(openPorts) == 0 {
		log.Println("No open ports found.")
		return
	}

	log.Println("Open ports:")
	for _, port := range openPorts {
		fmt.Printf(" - %d open\n", port)
	}
}

type data struct {
	host        string
	startPort   int
	endPort     int
	concurrency int
	timeout     int
}

func main() {
	info := data{}
	flag.StringVar(&info.host, "host", "localhost", "Target host to scan")
	flag.IntVar(&info.startPort, "start", 1, "Start port")
	flag.IntVar(&info.endPort, "end", 1024, "End port")
	flag.IntVar(&info.concurrency, "workers", 100, "Number of concurrent workers")
	flag.IntVar(&info.timeout, "timeout", 500, "Timeout in milliseconds per connection")
	flag.Parse()

	StartScanner(info)
}
