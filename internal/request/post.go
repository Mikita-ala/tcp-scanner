package request

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"tcp-scaner/internal/tcp"
	"time"
)

var tasks sync.Map

func StartScanner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var input ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		TaskID := uuid.New().String()
		slog.Info("Generate UUID" + TaskID)
		tasks.Store(TaskID, ScanResult{Status: ScanStatusPending})
		go func(id string, req ScanRequest) {
			scanner := tcp.Scanner{
				Address:     input.Host,
				StartPort:   input.StartPort,
				EndPort:     input.EndPort,
				Concurrency: input.Concurrency,
				Timeout:     time.Duration(input.Timeout) * time.Millisecond,
			}
			openPorts := scanner.Run()
			if openPorts != nil {
				tasks.Store(TaskID, ScanResult{Status: ScanStatusSuccess, Result: openPorts})
			}
			slog.Info("Finished scanner" + TaskID)
		}(TaskID, input)

		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{"task_id": TaskID})
	}
}

func StartScannerLocal(data ScanRequest) {
	scanner := tcp.Scanner{
		Address:     data.Host,
		StartPort:   data.StartPort,
		EndPort:     data.EndPort,
		Concurrency: data.Concurrency,
		Timeout:     time.Duration(data.Timeout) * time.Millisecond,
	}

	if data.StartPort < 1 || data.EndPort > 65535 || data.StartPort > data.EndPort {
		fmt.Println("Invalid port range")
		return
	}

	log.Printf("Scanning %s ports %d–%d with %d workers and timeout %dms...\n",
		data.Host, data.StartPort, data.EndPort, data.Concurrency, data.Timeout)

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
