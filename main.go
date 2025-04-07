package main

import (
	"log"
	"net/http"

	handlers "github.com/elitekentoy/chirpy/handlers/readiness"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(root))))
	serveMux.HandleFunc("/healthz", handlers.Readiness)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + listeningPort,
	}

	log.Printf("Serving files from %s on port: %s\n", root, listeningPort)
	log.Fatal(server.ListenAndServe())
}
