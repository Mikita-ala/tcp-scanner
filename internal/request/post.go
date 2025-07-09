package request

import (
	"encoding/json"
	"github.com/google/uuid"
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
		slog.Info("Generate UUID", TaskID)
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
			slog.Info("%s :Open ports on %s: %v\n", TaskID, input.Host, openPorts)
		}(TaskID, input)

		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{"task_id": TaskID})
	}
}
