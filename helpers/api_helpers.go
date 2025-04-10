package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/properties"
	"github.com/go-http-utils/headers"
)

func RespondToClient(writer http.ResponseWriter, payload any, status int) {
	writer.Header().Set(headers.ContentType, commons.APPLICATION_JSON)
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(payload); err != nil {
		RespondWithError(writer, properties.SERIALIZING_ISSUE, http.StatusInternalServerError)
	}
}

func RespondWithError(writer http.ResponseWriter, message string, status int) {
	http.Error(writer, message, status)
}
