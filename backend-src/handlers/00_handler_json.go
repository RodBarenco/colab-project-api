package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/RodBarenco/colab-project-api/rsakeys"
	"gorm.io/datatypes"
)

// ERROR RESPONSE
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
	w.WriteHeader(code)

	RespondWithJSON(w, code, errorResponse)
}

// JSON RESPONSE
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	if code < http.StatusOK || code >= http.StatusInternalServerError {
		code = http.StatusOK
	}

	data, err := json.Marshal(convertToJSON(payload))
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

// YOU CAN USE THIS FUNCTION IN CASES THAT YOU NEED TO RESPOND WITH ENCRYPTED DATA
func RespondWithEncryptedJSON(w http.ResponseWriter, code int, payload interface{}, publicKey string) {
	if code < http.StatusOK || code >= http.StatusInternalServerError {
		code = http.StatusOK
	}

	// Converter o payload para JSON
	data, err := json.Marshal(convertToJSON(payload))
	if err != nil {
		log.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Converter o payload para bytes
	payloadBytes := []byte(data)

	// Encriptar o payload
	encryptedData, err := rsakeys.EncryptAndEncode(publicKey, payloadBytes)
	if err != nil {
		log.Printf("Failed to encrypt payload: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to encrypt data")
		return
	}

	// Montar a resposta com o formato especificado
	response := map[string]string{
		"encrypted_data": encryptedData,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal encrypted response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataSize := len(responseJSON)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(dataSize))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(code)
	if _, err := w.Write(responseJSON); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}

// convertToJSON converts the given data into JSON format, preparing it for the server's JSON response.
func convertToJSON(data interface{}) interface{} {
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
