package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	// define chirp
	type chirp struct {
		Body string `json:"body"`
	}

	type error struct {
		Error string `json:"error"`
	}

	type cleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// define things
	decoder := json.NewDecoder(req.Body)
	data := chirp{}

	// Decode JSON
	err := decoder.Decode(&data)

	// if error occurs during decoding
	if err != nil {
		message := error{
			Error: "Something went wrong",
		}

		response, _ := json.Marshal(message)

		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(response)
		return
	}

	// if chirp is too long
	if len(data.Body) > chirpLimit {
		message := error{
			Error: "Chirp is too long",
		}

		response, _ := json.Marshal(message)

		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(response)
		return
	}

	message := cleaned{
		CleanedBody: removeProfanity(data.Body),
	}

	response, _ := json.Marshal(message)

	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

func removeProfanity(phrase string) string {
	blacklisted := []string{"kerfuffle", "sharbert", "fornax"}
	mask := "****"
	words := strings.Split(phrase, " ")

	for index, word := range words {
		if contains(blacklisted, word) {
			words[index] = mask
		}
	}

	return strings.Join(words, " ")
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if strings.ToLower(item) == strings.ToLower(target) {
			return true
		}
	}
	return false
}
