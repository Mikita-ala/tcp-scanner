package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tcp-scaner/internal/request"
)

func main() {
	r := mux.NewRouter()
	log.Println("Starting server")

	r.HandleFunc("/api/scan/{id}", request.GetScanResult()).Methods("GET")
	r.HandleFunc("/api/scan", request.StartScanner()).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
