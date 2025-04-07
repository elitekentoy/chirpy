package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		os.Exit(1)
	}

	dbQueries := database.New(db)

	apiConfig := &apiConfig{
		FileserverHits: atomic.Int32{},
		Database:       dbQueries,
	}

	serveMux := http.NewServeMux()

	// Setup file serving with middleware to count hits
	serveMux.Handle("/app/", apiConfig.MiddlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(root)))))

	// Define health check endpoint
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Define validate chirp endpoint
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

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
