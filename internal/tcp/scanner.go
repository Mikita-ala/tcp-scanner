package tcp

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

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
