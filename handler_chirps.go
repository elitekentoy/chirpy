package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpRequestBody struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (config *apiConfig) handlerChirps(writer http.ResponseWriter, req *http.Request) {

	request := ChirpRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(writer, "error decoding the request", http.StatusInternalServerError)
	}

	dbChirp, err := config.Database.CreateChirp(req.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      request.Body,
		UserID: uuid.NullUUID{
			UUID:  request.UserID,
			Valid: true,
		},
	})

	if err != nil {
		http.Error(writer, "error processing to the database", http.StatusInternalServerError)
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID.UUID,
	}

	data, err := json.Marshal(chirp)
	if err != nil {
		http.Error(writer, "error serializing chirp", http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")

	writer.Write(data)
}
