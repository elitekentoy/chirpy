package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot initialize database")
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

	// Define users endpoint
	serveMux.HandleFunc("POST /api/users", apiConfig.handlerUsers)

	// Define chirps endpoint
	serveMux.HandleFunc("POST /api/chirps", apiConfig.handlerChirps)

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
