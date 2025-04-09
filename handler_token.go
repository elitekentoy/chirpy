package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
)

func (config *apiConfig) handlerRefreshToken(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(writer, "invalid request", http.StatusBadRequest)
		return
	}

	dbToken, err := config.Database.GetRefreshTokenByToken(req.Context(), headerToken)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, "token not found", http.StatusBadRequest)
			return
		}

		http.Error(writer, "database communication error", http.StatusInternalServerError)
		return
	}

	token := models.RefreshTokenFromDatabase(dbToken)
	if token.IsExpired() {
		http.Error(writer, "token has expired", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(token.UserID, config.ApiSecret, time.Duration(ACCESS_TOKEN_EXPIRY_IN_HOURS)*time.Hour)
	if err != nil {
		http.Error(writer, "error creating token", http.StatusInternalServerError)
		return
	}

	content := map[string]string{}
	content["token"] = accessToken

	data, err := json.Marshal(content)
	if err != nil {
		http.Error(writer, "error serializing token", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	writer.Write(data)
}

func (config *apiConfig) handlerRevokeRefreshToken(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(writer, "invalid request", http.StatusBadRequest)
	}

	err = config.Database.UpdateRefreshTokenOnRevoke(req.Context(), database.UpdateRefreshTokenOnRevokeParams{
		Token: headerToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(writer, "token not found", http.StatusBadRequest)
		}

		http.Error(writer, "error in communication with database", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
