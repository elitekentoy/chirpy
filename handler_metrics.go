package main

import (
	"fmt"
	"net/http"

	"github.com/elitekentoy/chirpy/helpers"
)

func (config *apiConfig) handlerMetrics(writer http.ResponseWriter, req *http.Request) {
	template := `<html>
						<body>
							<h1>Welcome, Chirpy Admin</h1>
						<p>Chirpy has been visited %d times!</p>
						</body>
					</html>`

	helpers.RespondToClientWithHTMLBody(writer, fmt.Sprintf(template, config.FileserverHits.Load()), http.StatusOK)
}
