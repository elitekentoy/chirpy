package main

import (
	"fmt"
	"net/http"
)

func (config *apiConfig) handlerHits(writer http.ResponseWriter, req *http.Request) {
	// Set Content Type
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Set the status code
	writer.WriteHeader(http.StatusOK)

	data := fmt.Sprintf("Hits: %d", config.FileserverHits.Load())

	// Write the response body
	writer.Write([]byte(data))
}
