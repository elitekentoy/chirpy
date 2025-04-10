package models

import (
	"encoding/json"
	"net/http"

	"github.com/elitekentoy/chirpy/commons"
)

type PolkaRequestBody struct {
	Event string
	Data  map[string]string
}

func DecodePolkaRequestBody(req *http.Request) (PolkaRequestBody, error) {
	body := PolkaRequestBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)

	if err != nil {
		return PolkaRequestBody{}, err
	}

	return body, nil
}

func (request *PolkaRequestBody) IsUserUpgradeEvent() bool {
	return request.Event == commons.USER_UPGRADED
}

func (request *PolkaRequestBody) GetUserID() string {
	return request.Data["user_id"]
}
