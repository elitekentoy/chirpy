package main

import (
	"net/http"

	"github.com/elitekentoy/chirpy/commons"
	"github.com/elitekentoy/chirpy/helpers"
	"github.com/go-http-utils/headers"
)

func handlerReadiness(writer http.ResponseWriter, req *http.Request) {

	req.Header.Set(headers.ContentType, commons.TEXT_PLAIN)

	helpers.RespondToClientWithPlainBody(writer, "OK", http.StatusOK)
}
