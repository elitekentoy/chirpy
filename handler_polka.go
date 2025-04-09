package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/google/uuid"
)

type PolkaRequestBody struct {
	Event string
	Data  map[string]string
}

func (config *apiConfig) handlerPolka(writer http.ResponseWriter, req *http.Request) {

	apiKey, err := auth.GetAPIKey(req.Header)
	if apiKey == "" || err != nil || apiKey != config.PolkaSecret {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	request := PolkaRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&request)

	if err != nil {
		http.Error(writer, "error deserializing request", http.StatusInternalServerError)
		return
	}

	if request.Event != "user.upgraded" {
		http.Error(writer, "we dont care", http.StatusNoContent)
		return
	}

	userId := request.Data["user_id"]
	uid, err := uuid.Parse(userId)
	if err != nil {
		http.Error(writer, "error parsing user id", http.StatusInternalServerError)
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
		if err == sql.ErrNoRows {
			http.Error(writer, "user not found", http.StatusNotFound)
			return
		}

		http.Error(writer, "error communciting to the database", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
