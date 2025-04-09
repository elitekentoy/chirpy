package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return auth, errors.New("auth token does not exist")
	}

	token, found := strings.CutPrefix(auth, "Bearer ")
	if !found {
		return "", errors.New("invalid token")
	}

	return token, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	// Define Custom Claims Struct
	type CustomClaims struct {
		ID string
		jwt.RegisteredClaims
	}

	// Parse the JWT Token with Custom Claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	// If error occur during parsing
	if err != nil {
		log.Println(err)
		return uuid.Nil, errors.New("error parsing token")
	}

	// Assert the token claims as Custom Claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return uuid.Nil, errors.New("error getting claims")
	}

	// Parse the UUID from the "id" claim
	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("error parsing subject as uuid")
	}

	return uid, nil
}
