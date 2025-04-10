package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot initialize database")
	}

	apiConfig := &apiConfig{
		FileserverHits: atomic.Int32{},
		Database:       database.New(db),
		ApiSecret:      os.Getenv("SECRET_KEY"),
		PolkaSecret:    os.Getenv("POLKA_SECRET_KEY"),
	}

	serveMux := http.NewServeMux()

	// Setup file serving with middleware to count hits
	serveMux.Handle("/app/", apiConfig.MiddlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(commons.ROOT)))))

	// Define health check endpoint
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Define get one chirp by ID
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirp)

	// Define get all chirps endpoint
	serveMux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)

	// Define validate chirp endpoint
	serveMux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	// Define users endpoint
	serveMux.HandleFunc("POST /api/users", apiConfig.handlerRegisterUser)

	// Define create chirp endpoint
	serveMux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)

	// Define login endpoing
	serveMux.HandleFunc("POST /api/login", apiConfig.handlerLogin)

	// Define refresh token endpoint
	serveMux.HandleFunc("POST /api/refresh", apiConfig.handlerRefreshToken)

	// Define revoke token endpoint
	serveMux.HandleFunc("POST /api/revoke", apiConfig.handlerRevokeRefreshToken)

	// Define polka webhooks
	serveMux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlerPolka)

	// Define update user endpoint
	serveMux.HandleFunc("PUT /api/users", apiConfig.handlerUpdateUserDetails)

	// Define delete chip endpoint
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.handlerDeleteChirp)

	// Define metric endpoint
	serveMux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)

	// Define reset endpoint
	serveMux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)

	// Setup HTTP Server
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + commons.LISTENING_PORT,
	}

	log.Printf("Serving files from %s on port: %s\n", commons.ROOT, commons.LISTENING_PORT)
	log.Fatal(server.ListenAndServe())
}
