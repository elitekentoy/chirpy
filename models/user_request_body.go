package models

import (
	"encoding/json"
	"net/http"
)

type UsersRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func DecodeUsersRequestBody(req *http.Request) (UsersRequestBody, error) {
	var body UsersRequestBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return UsersRequestBody{}, err
	}
	return body, nil
}

func (u UsersRequestBody) IsInvalid() bool {
	return u.Email == "" || u.Password == ""
}
