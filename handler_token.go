package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/helpers"
	"github.com/elitekentoy/chirpy/internal/auth"
	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/elitekentoy/chirpy/models"
	"github.com/elitekentoy/chirpy/properties"
)

func (config *apiConfig) handlerRefreshToken(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	dbToken, err := config.Database.GetRefreshTokenByToken(req.Context(), headerToken)
	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	token := models.RefreshTokenFromDatabase(dbToken)
	if token.IsExpired() {
		helpers.RespondWithError(writer, properties.EXPIRED_TOKEN, http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.MakeJWT(token.UserID, config.ApiSecret, commons.ACCESS_TOKEN_EXPIRY)
	if err != nil {
		helpers.RespondWithError(writer, properties.TOKEN_GENERIC_ERROR, http.StatusInternalServerError)
		return
	}

	content := models.TokenResponse{
		Token: accessToken,
	}

	helpers.RespondToClient(writer, content, http.StatusOK)
}

func (config *apiConfig) handlerRevokeRefreshToken(writer http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		helpers.RespondWithError(writer, properties.INVALID_TOKEN, http.StatusUnauthorized)
		return
	}

	err = config.Database.UpdateRefreshTokenOnRevoke(req.Context(), database.UpdateRefreshTokenOnRevokeParams{
		Token: headerToken,
		RevokedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		helpers.HandleDatabaseError(writer, err)
		return
	}

	helpers.RespondToClient(writer, nil, http.StatusNoContent)
}
