package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
	"github.com/google/uuid"
)

type UsersRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (config *apiConfig) handlerUsers(writer http.ResponseWriter, req *http.Request) {

	request := UsersRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(writer, "error occured in decoding", http.StatusInternalServerError)
		return
	}

	if request.Email == "" || request.Password == "" {
		http.Error(writer, "email or password cannot be empty", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		http.Error(writer, "cannot hash password", http.StatusInternalServerError)
		return
	}

	dbUser, err := config.Database.CreateUser(req.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})

	user := models.UserFromDatabase(dbUser)

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
