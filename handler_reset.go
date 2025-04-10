package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/helpers"
	"github.com/elitekentoy/chirpy/properties"
)

func (config *apiConfig) handlerReset(writer http.ResponseWriter, req *http.Request) {

	if os.Getenv(commons.PLATFORM_KEY) != commons.DEV {
		helpers.RespondWithError(writer, properties.NO_PERMISSIONS, http.StatusForbidden)
		return
	}

	err := config.Database.DeleteUsers(req.Context())
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	config.FileserverHits.Store(commons.DEFAULT_HITS)

	helpers.RespondToClientWithBody(writer, fmt.Sprintf("Hits: %d", config.FileserverHits.Load()), http.StatusOK)
}
