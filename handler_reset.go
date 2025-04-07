package main

import (
	"fmt"
	"net/http"
	"os"
)

func (config *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {

	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		http.Error(writer, "you do not have permissions for reset", http.StatusForbidden)
		return
	}

	err := config.Database.DeleteUsers(req.Context())
	if err != nil {
		http.Error(writer, "cannot perform operations in database", http.StatusInternalServerError)
		return
	}

	config.FileserverHits.Store(0)

	// Set Content Type
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Set the status code
	writer.WriteHeader(http.StatusOK)

	data := fmt.Sprintf("Hits: %d", config.FileserverHits.Load())

	// Write the response body
	writer.Write([]byte(data))
}
