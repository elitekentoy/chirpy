package handlers

import "net/http"

func Readiness(writer http.ResponseWriter, req *http.Request) {
	// Set Content Type
	req.Header.Set("content-type", "text/plain; charset=utf-8")

	// Set the status coke
	writer.WriteHeader(http.StatusOK)

	// Write the response body
	writer.Write([]byte("OK"))
}
