package main

import (
	"database/sql"
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

func (config *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {

	request := UsersRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)

	if err != nil {
		http.Error(writer, "error deserializing request", http.StatusInternalServerError)
		return
	}

	if request.Email == "" || request.Password == "" {
		http.Error(writer, "email or password cannot be empty", http.StatusBadRequest)
		return
	}

	dbUser, err := config.Database.GetUserByEmail(req.Context(), request.Email)
	if err != nil {

		if err == sql.ErrNoRows {
			http.Error(writer, "user not found", http.StatusNotFound)
			return
		}

		http.Error(writer, "error occured in communicating to the database", http.StatusInternalServerError)
		return
	}

	valid := auth.CheckPasswordHash(dbUser.HashedPassword, request.Password) == nil

	if !valid {
		http.Error(writer, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, config.ApiSecret, time.Duration(ACCESS_TOKEN_EXPIRY_IN_HOURS)*time.Hour)
	if err != nil {
		http.Error(writer, "cannot create token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		http.Error(writer, "cannot create refresh token", http.StatusInternalServerError)
	}

	dbRefreshToken, err := config.Database.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: uuid.NullUUID{
			UUID:  dbUser.ID,
			Valid: true,
		},
		ExpiresAt: time.Now().Add(time.Duration(REFRESH_TOKEN_EXPIRY_IN_HOURS) * time.Hour),
	})

	if err != nil {
		http.Error(writer, "error creating refresh token", http.StatusInternalServerError)
		return
	}

	user := models.UserFromDatabaseWithTokens(dbUser, accessToken, dbRefreshToken.Token)
	data, err := json.Marshal(user)

	if err != nil {
		http.Error(writer, "error in serializing user", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	writer.Write(data)
}

func (config *apiConfig) handlerUpdateUserDetails(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil || headerToken == "" {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	uid, err := auth.ValidateJWT(headerToken, config.ApiSecret)
	if uid == uuid.Nil || err != nil {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	request := UsersRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&request)

	if err != nil {
		http.Error(writer, "error in deserializing request", http.StatusInternalServerError)
		return
	}

	if request.Email == "" || request.Password == "" {
		http.Error(writer, "email or password cannot be empty", http.StatusBadRequest)
		return
	}

	dbUser, err := config.Database.GetUserByID(req.Context(), uid)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, "user not found", http.StatusBadRequest)
			return
		}

		http.Error(writer, "error in communicating to database", http.StatusInternalServerError)
		return
	}

	if uid != dbUser.ID {
		http.Error(writer, "invalid token", http.StatusUnauthorized)
		return
	}

	hashed, err := auth.HashPassword(request.Password)
	if err != nil {
		http.Error(writer, "error hashing password", http.StatusInternalServerError)
		return
	}

	dbUser, err = config.Database.UpdateUser(req.Context(), database.UpdateUserParams{
		UpdatedAt:      time.Now(),
		Email:          request.Email,
		HashedPassword: hashed,
		ID:             uid,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, "user not found", http.StatusBadRequest)
			return
		}

		http.Error(writer, "error in communicating to database", http.StatusInternalServerError)
		return
	}

	user := models.UserFromDatabase(dbUser)

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(writer, "error serializing user", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	writer.Write(data)
}
