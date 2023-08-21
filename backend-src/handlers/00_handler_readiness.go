package handlers

import (
	"net/http"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, gormJSON{})
}

func HandlerReadinessEncrypted(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	RespondWithJSON(w, http.StatusOK, gormJSON{})
}

type gormJSON struct{}
