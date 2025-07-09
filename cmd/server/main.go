package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"tcp-scaner/internal/request"
)

func main() {
	var port string
	flag.StringVar(&port, "port", ":8000", "Target port to server")

	r := mux.NewRouter()
	log.Println("Starting server, port:", port)
	log.Printf("http://localhost%s", port)

	r.HandleFunc("/api/scan/{id}", request.GetScanResult()).Methods("GET")
	r.HandleFunc("/api/scan", request.StartScanner()).Methods("POST")

	log.Fatal(http.ListenAndServe(port, r))
}
