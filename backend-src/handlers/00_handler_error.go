package handlers

import (
	"net/http"
)

func HandlerError(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusBadRequest, "Something went wrong")
}

type gormError struct {
	Error string `json:"error"`
}
