package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/google/uuid"
)

type requestBody struct {
	Email string `json:"email"`
}

func (config *apiConfig) handlerUsers(writer http.ResponseWriter, req *http.Request) {

	request := requestBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(writer, "error occured in decoding", http.StatusInternalServerError)
		return
	}

	user, err := config.Database.CreateUser(req.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     request.Email,
	})

	if err != nil {
		http.Error(writer, "error occured in database", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)

	if err != nil {
		http.Error(writer, "failed to serialize user", http.StatusInternalServerError)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	writer.Write(data)
}
