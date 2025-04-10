package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/elitekentoy/chirpy/helpers"
	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
	"github.com/elitekentoy/chirpy/properties"
	"github.com/google/uuid"
)

func (config *apiConfig) handlerCreateChirp(writer http.ResponseWriter, req *http.Request) {

	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(authToken, config.ApiSecret)
	if err != nil {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	request, err := models.DecodeChirpRequestBody(req)
	if err != nil {
		helpers.RespondWithError(writer, properties.DESERIALIZING_ISSUE, http.StatusInternalServerError)
		return
	}

	dbChirp, err := config.Database.CreateChirp(req.Context(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body:      request.Body,
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
	})

	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	helpers.RespondToClient(writer, models.ChirpFromDatabase(dbChirp), http.StatusCreated)
}

func (config *apiConfig) handlerGetChirps(writer http.ResponseWriter, req *http.Request) {
	authorID := req.URL.Query().Get("author_id")
	uid, err := uuid.Parse(authorID)
	if err != nil && authorID != "" {
		helpers.RespondWithError(writer, properties.PARSING_ISSUE, http.StatusInternalServerError)
		return
	}

	dbChirps := []database.Chirp{}

	if authorID == "" {
		dbChirps, err = config.Database.GetChirps(req.Context())
	} else {
		dbChirps, err = config.Database.GetChirpsByUserID(req.Context(), uuid.NullUUID{
			UUID:  uid,
			Valid: true,
		})
	}

	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	if req.URL.Query().Get("sort") == "desc" {
		sort.Slice(dbChirps, func(i, j int) bool {
			return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
		})
	}

	chirps := []models.Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, models.ChirpFromDatabase(chirp))
	}

	helpers.RespondToClient(writer, chirps, http.StatusOK)
}

func (config *apiConfig) handlerGetChirp(writer http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")

	if chirpID == "" {
		http.Error(writer, properties.MISSING_ID, http.StatusBadRequest)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(writer, properties.INCORRECT_INPUT, http.StatusBadRequest)
		return
	}

	dbChirp, err := config.Database.GetChirpByID(req.Context(), chirpUUID)
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	helpers.RespondToClient(writer, models.ChirpFromDatabase(dbChirp), http.StatusOK)
}

func (config *apiConfig) handlerDeleteChirp(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if headerToken == "" || err != nil {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	uid, err := auth.ValidateJWT(headerToken, config.ApiSecret)
	if uid == uuid.Nil || err != nil {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	chirpID := req.PathValue("chirpID")

	if chirpID == "" {
		helpers.RespondWithError(writer, properties.MISSING_ID, http.StatusBadRequest)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		helpers.RespondWithError(writer, properties.PARSING_ISSUE, http.StatusInternalServerError)
		return
	}

	dbChirp, err := config.Database.GetChirpByID(req.Context(), chirpUUID)
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	chirp := models.ChirpFromDatabase(dbChirp)
	if uid != chirp.UserID {
		helpers.RespondWithError(writer, properties.NO_PERMISSIONS, http.StatusForbidden)
		return
	}

	err = config.Database.DeleteChirpByID(req.Context(), chirp.ID)
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	helpers.RespondToClient(writer, nil, http.StatusNoContent)
}
