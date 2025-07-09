package main

import (
	"flag"
	"tcp-scaner/internal/request"
)

func main() {
	info := request.ScanRequest{}

	flag.StringVar(&info.Host, "host", "localhost", "Target host to scan")
	flag.IntVar(&info.StartPort, "start", 1, "Start port")
	flag.IntVar(&info.EndPort, "end", 1024, "End port")
	flag.IntVar(&info.Concurrency, "workers", 100, "Number of concurrent workers")
	flag.IntVar(&info.Timeout, "timeout", 500, "Timeout in milliseconds per connection")
	flag.Parse()

	request.StartScannerLocal(info)
}
