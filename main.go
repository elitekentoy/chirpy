package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	apiConfig := &apiConfig{
		FileserverHits: atomic.Int32{},
	}

	serveMux := http.NewServeMux()

	// Setup file serving with middleware to count hits
	serveMux.Handle("/app/", apiConfig.MiddlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

	// Define health check endpoint
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Define metric endpoint
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)

	// Define reset endpoint
	serveMux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)

	// Setup HTTP Server
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + listeningPort,
	}

	log.Printf("Serving files from %s on port: %s\n", root, listeningPort)
	log.Fatal(server.ListenAndServe())
}
