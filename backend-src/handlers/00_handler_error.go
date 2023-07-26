package handlers

import (
	"net/http"
)

func HandlerError(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusUnauthorized, "Something went wrong")
}
