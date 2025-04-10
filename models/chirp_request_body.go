package models

import (
	"encoding/json"
	"net/http"

	"github.com/elitekentoy/chirpy/commons"
)

type ChirpRequestBody struct {
	Body string `json:"body"`
}

func DecodeChirpRequestBody(req *http.Request) (ChirpRequestBody, error) {
	body := ChirpRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)

	if err != nil {
		return ChirpRequestBody{}, err
	}

	return body, nil
}

func (chirp *ChirpRequestBody) ExceedsLimit() bool {
	return len(chirp.Body) > commons.CHIRP_LIMIT
}
