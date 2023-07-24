package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/datatypes"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {

	if code < http.StatusBadRequest || code >= http.StatusInternalServerError {
		code = http.StatusInternalServerError
	}

	var friendlyMsg string
	switch code {
	case http.StatusBadRequest:
		friendlyMsg = "Bad Request"
	case http.StatusUnauthorized:
		friendlyMsg = "Unauthorized"
	case http.StatusForbidden:
		friendlyMsg = "Forbidden"
	case http.StatusNotFound:
		friendlyMsg = "Not Found"
	case http.StatusInternalServerError:
		friendlyMsg = "Internal Server Error"
	default:
		friendlyMsg = "Error"
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	errorResponse := ErrorResponse{
		Error: friendlyMsg + ": " + msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	RespondWithJSON(w, code, errorResponse)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	if code < http.StatusOK || code >= http.StatusInternalServerError {
		code = http.StatusOK
	}

	data, err := json.Marshal(toGORMJSON(payload))
	if err != nil {
		log.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataSize := len(data)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(dataSize))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

func toGORMJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case datatypes.JSON:
		var result interface{}
		if err := v.UnmarshalJSON([]byte(v)); err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
		}
		return result
	case []datatypes.JSON:
		var result []interface{}
		for _, d := range v {
			var r interface{}
			if err := d.UnmarshalJSON([]byte(d)); err != nil {
				log.Printf("Failed to unmarshal JSON: %v", err)
				continue
			}
			result = append(result, r)
		}
		return result
	default:
		return data
	}
}
