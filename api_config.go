package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	FileserverHits atomic.Int32
}

func (config *apiConfig) MiddlewareMetricInc(next http.Handler) http.Handler {

	// Call the next handler in the chain
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware called, incrementing hit counter")
		config.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
