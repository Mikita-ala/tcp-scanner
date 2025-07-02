package tcp

import (
	"time"
)

type Scanner struct {
	Address     string
	StartPort   int
	EndPort     int
	Concurrency int
	Timeout     time.Duration
}
