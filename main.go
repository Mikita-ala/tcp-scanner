package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

type Scanner struct {
	Address     string
	StartPort   int
	EndPort     int
	Concurrency int
	Timeout     time.Duration
}

func (s *Scanner) ScanPort(port int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	target := fmt.Sprintf("%s:%d", s.Address, port)
	conn, err := net.DialTimeout("tcp", target, s.Timeout)
	if err == nil {
		conn.Close()
		results <- port
	}
}

func (s *Scanner) Run() []int {
	var wg sync.WaitGroup
	ports := make(chan int, s.Concurrency)
	results := make(chan int, s.EndPort-s.StartPort+1)

	// запускаем воркеры
	for i := 0; i < s.Concurrency; i++ {
		go func() {
			for port := range ports {
				s.ScanPort(port, results, &wg)
			}
		}()
	}

	// добавляем порты в очередь
	for port := s.StartPort; port <= s.EndPort; port++ {
		wg.Add(1)
		ports <- port
	}

	wg.Wait()
	close(ports)
	close(results)

	var openPorts []int
	for port := range results {
		openPorts = append(openPorts, port)
	}
	sort.Ints(openPorts)
	return openPorts
}

func main() {
	var (
		host        string
		startPort   int
		endPort     int
		concurrency int
		timeout     int
	)

	flag.StringVar(&host, "host", "localhost", "Target host to scan")
	flag.IntVar(&startPort, "start", 1, "Start port")
	flag.IntVar(&endPort, "end", 1024, "End port")
	flag.IntVar(&concurrency, "workers", 100, "Number of concurrent workers")
	flag.IntVar(&timeout, "timeout", 500, "Timeout in milliseconds per connection")
	flag.Parse()

	if startPort < 1 || endPort > 65535 || startPort > endPort {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}

	scanner := Scanner{
		Address:     host,
		StartPort:   startPort,
		EndPort:     endPort,
		Concurrency: concurrency,
		Timeout:     time.Duration(timeout) * time.Millisecond,
	}

	log.Printf("Scanning %s ports %d–%d with %d workers and timeout %dms...\n",
		host, startPort, endPort, concurrency, timeout)

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
