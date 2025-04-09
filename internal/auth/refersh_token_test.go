package auth

import (
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil || token == "" {
		t.Error("error making refresh toking")
	}
}
