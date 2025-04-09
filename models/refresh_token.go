package models

import (
	"time"

	"github.com/elitekentoy/chirpy/internal/database"
	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string     `json:"token"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UserID    uuid.UUID  `json:"user_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

func RefreshTokenFromDatabase(dbToken database.RefreshToken) RefreshToken {

	var token *time.Time
	if dbToken.RevokedAt.Valid {
		token = &dbToken.RevokedAt.Time
	}

	return RefreshToken{
		Token:     dbToken.Token,
		CreatedAt: dbToken.CreatedAt,
		UpdatedAt: dbToken.UpdatedAt,
		UserID:    dbToken.UserID.UUID,
		ExpiresAt: dbToken.ExpiresAt,
		RevokedAt: token,
	}
}

func (token *RefreshToken) IsExpired() bool {
	return time.Now().After(token.ExpiresAt) || token.RevokedAt != nil
}
