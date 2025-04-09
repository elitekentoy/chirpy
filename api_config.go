package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/elitekentoy/chirpy/internal/database"
)

type apiConfig struct {
	FileserverHits atomic.Int32
	Database       *database.Queries
	ApiSecret      string
}

func (config *apiConfig) MiddlewareMetricInc(next http.Handler) http.Handler {

	// Call the next handler in the chain
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware called, incrementing hit counter")
		config.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
