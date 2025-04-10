package main

import (
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/helpers"
	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
	"github.com/elitekentoy/chirpy/properties"
	"github.com/google/uuid"
)

func (config *apiConfig) handlerRegisterUser(writer http.ResponseWriter, req *http.Request) {

	request, err := models.DecodeUsersRequestBody(req)
	if err != nil {
		helpers.RespondWithError(writer, properties.DESERIALIZING_ISSUE, http.StatusBadRequest)
		return
	}

	if request.IsInvalid() {
		helpers.RespondWithError(writer, properties.LACKING_EMAIL_OR_PASSWORD, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		helpers.RespondWithError(writer, properties.HASHING_PASSWORD_ISSUE, http.StatusInternalServerError)
		return
	}

	dbUser, err := config.Database.CreateUser(req.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		helpers.RespondWithError(writer, properties.GENERIC_DATABASE_ERROR, http.StatusInternalServerError)
		return
	}

	helpers.RespondToClient(writer, models.UserFromDatabase(dbUser), http.StatusCreated)
}

func (config *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {

	request, err := models.DecodeUsersRequestBody(req)

	if err != nil {
		helpers.RespondWithError(writer, properties.DESERIALIZING_ISSUE, http.StatusBadRequest)
		return
	}

	if request.IsInvalid() {
		helpers.RespondWithError(writer, properties.LACKING_EMAIL_OR_PASSWORD, http.StatusBadRequest)
		return
	}

	dbUser, err := config.Database.GetUserByEmail(req.Context(), request.Email)
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	if err := auth.CheckPasswordHash(dbUser.HashedPassword, request.Password); err != nil {
		helpers.RespondWithError(writer, properties.INCORRECT_INPUT, http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(dbUser.ID, config.ApiSecret, commons.ACCESS_TOKEN_EXPIRY)
	if err != nil {
		helpers.RespondWithError(writer, properties.TOKEN_GENERIC_ERROR, http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		helpers.RespondWithError(writer, properties.TOKEN_GENERIC_ERROR, http.StatusInternalServerError)
		return
	}

	dbRefreshToken, err := config.Database.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: uuid.NullUUID{
			UUID:  dbUser.ID,
			Valid: true,
		},
		ExpiresAt: time.Now().UTC().Add(commons.REFRESH_TOKEN_EXPIRY),
	})

	if err != nil {
		helpers.RespondWithError(writer, properties.TOKEN_GENERIC_ERROR, http.StatusInternalServerError)
		return
	}

	helpers.RespondToClient(writer, models.UserFromDatabaseWithTokens(dbUser, accessToken, dbRefreshToken.Token), http.StatusOK)
}

func (config *apiConfig) handlerUpdateUserDetails(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil || headerToken == "" {
		http.Error(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	uid, err := auth.ValidateJWT(headerToken, config.ApiSecret)
	if uid == uuid.Nil || err != nil {
		http.Error(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	request, err := models.DecodeUsersRequestBody(req)
	if err != nil {
		helpers.RespondWithError(writer, properties.DESERIALIZING_ISSUE, http.StatusBadRequest)
		return
	}

	if request.IsInvalid() {
		helpers.RespondWithError(writer, properties.LACKING_EMAIL_OR_PASSWORD, http.StatusBadRequest)
		return
	}

	dbUser, err := config.Database.GetUserByID(req.Context(), uid)
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	if uid != dbUser.ID {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	hashed, err := auth.HashPassword(request.Password)
	if err != nil {
		helpers.RespondWithError(writer, properties.HASHING_PASSWORD_ISSUE, http.StatusInternalServerError)
		return
	}

	updatedUser, err := config.Database.UpdateUser(req.Context(), database.UpdateUserParams{
		UpdatedAt:      time.Now().UTC(),
		Email:          request.Email,
		HashedPassword: hashed,
		ID:             uid,
	})

	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	helpers.RespondToClient(writer, models.UserFromDatabase(updatedUser), http.StatusOK)
}
