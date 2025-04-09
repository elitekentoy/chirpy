package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJwtValid(t *testing.T) {
	id, _ := uuid.Parse("cceeb45d-547c-41b6-8d11-d82836ce3bab")
	secret := "secret"
	expiresIn := time.Duration(10) * time.Second

	token, err := MakeJWT(id, secret, expiresIn)
	if token == "" || err != nil {
		t.Errorf("error testing makejwt() -> token :%s -> err :%v", token, err)
	}
}

func TestValidateJWTValid(t *testing.T) {
	id, _ := uuid.Parse("cceeb45d-547c-41b6-8d11-d82836ce3bab")
	secret := "secret"
	expiresIn := time.Duration(10) * time.Second

	token, _ := MakeJWT(id, secret, expiresIn)

	uid, err := ValidateJWT(token, secret)
	if uid == uuid.Nil || err != nil {
		t.Errorf("error testing validatejwt() -> uid :%v -> err :%v", uid, err)
	}
}

func TestValidateJWTInvalidSecret(t *testing.T) {
	id, _ := uuid.Parse("cceeb45d-547c-41b6-8d11-d82836ce3bab")
	secret := "secret"
	expiresIn := time.Duration(10) * time.Second

	token, _ := MakeJWT(id, secret, expiresIn)

	uid, err := ValidateJWT(token, "new_secret")
	if uid != uuid.Nil || err == nil {
		t.Errorf("error testing validatejwt() -> uid :%v -> err :%v", uid, err)
	}
}

func TestValidateJWTExpired(t *testing.T) {
	id, _ := uuid.Parse("cceeb45d-547c-41b6-8d11-d82836ce3bab")
	secret := "secret"
	expiresIn := time.Duration(-10) * time.Second

	token, _ := MakeJWT(id, secret, expiresIn)

	uid, err := ValidateJWT(token, secret)
	if uid != uuid.Nil || err == nil {
		t.Errorf("error testing validatejwt() -> uid :%v -> err :%v", uid, err)
	}
}

func TestGetBearerTokenValid(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer my_token_should_be_here")

	token, err := GetBearerToken(headers)

	if token == "" || err != nil {
		t.Errorf("seems like token is not found")
	}
}

func TestGetBearerTokenInvalid(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "my_token_should_be_here")

	token, err := GetBearerToken(headers)

	if token != "" || err == nil {
		t.Errorf("seems like token is not found")
	}
}

func TestGetBearerTokenHeaderNotExists(t *testing.T) {
	token, err := GetBearerToken(http.Header{})

	if token != "" || err == nil {
		t.Errorf("seems like token is not found")
	}
}
