package main

import (
	"database/sql"
	"net/http"

	"github.com/elitekentoy/chirpy/helpers"
	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
	"github.com/elitekentoy/chirpy/properties"
	"github.com/google/uuid"
)

func (config *apiConfig) handlerPolka(writer http.ResponseWriter, req *http.Request) {

	apiKey, err := auth.GetAPIKey(req.Header)
	if apiKey == "" || err != nil || apiKey != config.PolkaSecret {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	request, err := models.DecodePolkaRequestBody(req)
	if err != nil {
		helpers.RespondToClientWithBody(writer, properties.DESERIALIZING_ISSUE, http.StatusBadRequest)
		return
	}

	if !request.IsUserUpgradeEvent() {
		helpers.RespondWithError(writer, properties.CANNOT_PROCESS, http.StatusNoContent)
		return
	}

	uid, err := uuid.Parse(request.GetUserID())
	if err != nil {
		helpers.RespondWithError(writer, properties.INCORRECT_INPUT, http.StatusBadRequest)
		return
	}

	err = config.Database.UpdateChirpyRed(req.Context(), database.UpdateChirpyRedParams{
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ID: uid,
	})

	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	helpers.RespondToClient(writer, nil, http.StatusNoContent)
}
