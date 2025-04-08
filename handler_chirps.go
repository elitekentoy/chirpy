package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
	"github.com/google/uuid"
)

type ChirpRequestBody struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (config *apiConfig) handlerCreateChirp(writer http.ResponseWriter, req *http.Request) {

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

	chirp := models.ChirpFromDatabase(dbChirp)

	data, err := json.Marshal(chirp)
	if err != nil {
		http.Error(writer, "error serializing chirp", http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")

	writer.Write(data)
}

func (config *apiConfig) handlerGetChirps(writer http.ResponseWriter, req *http.Request) {
	dbChirps, err := config.Database.GetChirps(req.Context())
	if err != nil {
		http.Error(writer, "error connecting to the database", http.StatusInternalServerError)
	}

	chirps := []models.Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, models.ChirpFromDatabase(chirp))
	}

	data, err := json.Marshal(chirps)
	if err != nil {
		http.Error(writer, "error seraializing chirps", http.StatusInternalServerError)
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")

	writer.Write(data)
}

func (config *apiConfig) handlerGetChirp(writer http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

	if chirpID == "" {
		http.Error(writer, "chirp id is not specified", http.StatusBadRequest)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(writer, "error parsing chirp id", http.StatusInternalServerError)
		return
	}

	dbChirp, err := config.Database.GetChirpByID(req.Context(), chirpUUID)
	if err != nil {

		if err == sql.ErrNoRows {
			http.Error(writer, "chirp id is not found", http.StatusNotFound)
			return
		}

		http.Error(writer, "database error occured", http.StatusInternalServerError)
		return
	}

	chirp := models.ChirpFromDatabase(dbChirp)

	data, err := json.Marshal(chirp)
	if err != nil {
		http.Error(writer, "error in serializing chirp", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	writer.Write(data)
	return
}
