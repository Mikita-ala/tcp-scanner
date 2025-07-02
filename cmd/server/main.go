package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tcp-scaner/internal/tcp"
	"time"
)

type data struct {
	Host        string `json:"host"`
	StartPort   int    `json:"startPort"`
	EndPort     int    `json:"endPort"`
	Concurrency int    `json:"concurrency"`
	Timeout     int    `json:"timeout"`
}

func StartScanner() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var input data
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//
		//TaskID := uuid.New().String()
		//task

		scanner := tcp.Scanner{
			Address:     input.Host,
			StartPort:   input.StartPort,
			EndPort:     input.EndPort,
			Concurrency: input.Concurrency,
			Timeout:     time.Duration(input.Timeout) * time.Millisecond,
		}

		openPorts := scanner.Run()

		// Логируем в консоль
		fmt.Printf("Open ports on %s: %v\n", input.Host, openPorts)

		// Отправляем результат
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(openPorts); err != nil {
			http.Error(w, "failed to encode result", http.StatusInternalServerError)
		}
	}
}

func main() {
	r := mux.NewRouter()
	log.Println("Starting server")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	r.HandleFunc("/api/scan", StartScanner()).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
