package helpers

import (
	"database/sql"
	"net/http"

	"github.com/elitekentoy/chirpy/properties"
)

func HandleDatabaseError(writer http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		RespondWithError(writer, properties.RECORD_NOT_FOUND, http.StatusNotFound)
		return
	}

	RespondWithError(writer, properties.GENERIC_DATABASE_ERROR, http.StatusInternalServerError)
}
