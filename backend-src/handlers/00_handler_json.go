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

	// Convert the payload to JSON
	data, err := json.Marshal(convertToJSON(payload))
	if err != nil {
		log.Printf("Failed to marshal JSON response %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Convert the payload to bytes
	payloadBytes := []byte(data)

	var encryptedData string
	var aesKey []byte

	if len(payloadBytes) > 512 { // Use your desired threshold here
		// Generate AES key for large files
		aesKey, err = rsakeys.GenerateAESKeyForLargeFile()
		if err != nil {
			log.Printf("Failed to generate AES key: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to encrypt data")
			return
		}

		// Encrypt payload with AES
		encryptedData, err = rsakeys.EncryptAES(aesKey, payloadBytes)
		if err != nil {
			log.Printf("Failed to encrypt payload with AES: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to encrypt data")
			return
		}
	} else {
		// Encrypt payload with RSA and encode as base64
		encryptedData, err = rsakeys.EncryptAndEncode(publicKey, payloadBytes)
		if err != nil {
			log.Printf("Failed to encrypt payload: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to encrypt data")
			return
		}
	}

	response := map[string]string{
		"encrypted_data": encryptedData,
	}

	if aesKey != nil {
		rsaEncryptedAESKey, err := rsakeys.EncryptAndEncode(publicKey, aesKey)
		if err != nil {
			log.Printf("Failed to encrypt AES key: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Failed to encrypt data")
			return
		}
		response["aes_key"] = rsaEncryptedAESKey
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
