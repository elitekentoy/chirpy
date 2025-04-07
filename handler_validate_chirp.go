package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	// define chirp
	type chirp struct {
		Body string `json:"body"`
	}

	type error struct {
		Error string `json:"error"`
	}

	type valid struct {
		Valid bool `json:"valid"`
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

	message := valid{
		Valid: true,
	}

	response, _ := json.Marshal(message)

	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}
