package main

import (
	"net/http"
	"strings"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/helpers"
	"github.com/elitekentoy/chirpy/models"
	"github.com/elitekentoy/chirpy/properties"
)

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	request, err := models.DecodeChirpRequestBody(req)
	if err != nil {
		helpers.RespondToClientWithPlainBody(writer, models.ErrorResponse{Error: properties.GENERIC_ERROR}, http.StatusBadRequest)
		return
	}

	if request.ExceedsLimit() {
		helpers.RespondToClient(writer, models.ErrorResponse{Error: properties.CHIRP_EXCEEDED}, http.StatusBadRequest)
		return
	}

	cleaned := models.ChirpCleanedBody{CleanedBody: removeProfanity(request.Body)}

	helpers.RespondToClient(writer, cleaned, http.StatusOK)
}

func removeProfanity(phrase string) string {
	blacklisted := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(phrase, " ")

	for index, word := range words {
		if contains(blacklisted, word) {
			words[index] = commons.MASK
		}
	}

	return strings.Join(words, " ")
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, target) {
			return true
		}
	}
	return false
}
