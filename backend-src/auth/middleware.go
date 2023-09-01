package auth

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"

	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/RodBarenco/colab-project-api/rsakeys"
)

func AuthMiddleware(requiredRole string, secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the request header
		tokenString := extractTokenFromHeader(r)

		// If the token is not present in the header, try getting it from the request context
		if tokenString == "" {
			token, ok := r.Context().Value("jwtToken").(string)
			if !ok {
				http.Error(w, "JWT token not found in the request", http.StatusUnauthorized)
				return
			}
			tokenString = token
		}

		// Call the Authorize function for token validation
		statusCode, err := Authorize(tokenString, secret, requiredRole)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// Call the next handler with the updated context.
		next.ServeHTTP(w, r)
	})
}

func ActionsMiddleware(id, secret string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the request header
		tokenString := extractTokenFromHeader(r)

		// If the token is not present in the header, try getting it from the request context
		if tokenString == "" {
			token, ok := r.Context().Value("jwtToken").(string)
			if !ok {
				http.Error(w, "JWT token not found in the request", http.StatusUnauthorized)
				return
			}
			tokenString = token
		}

		// Call the Authorize function for token validation
		statusCode, err := AuthorizeActions(tokenString, secret, id)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// If the authorization is successful, call the handler function
		handler(w, r)
	}
}

func DecryptionMiddleware(next http.Handler, privateKey *rsa.PrivateKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the response is encrypted
		encryptedHeader := r.Header.Get("Encrypted")
		if encryptedHeader != "true" {
			next.ServeHTTP(w, r)
			return
		}

		// Read the response body
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Parse the JSON to extract the "encrypted_data" field
		var encryptedJSON map[string]string
		err = json.Unmarshal(buf.Bytes(), &encryptedJSON)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
			return
		}

		encryptedData, ok := encryptedJSON["encrypted_data"]
		if !ok {
			http.Error(w, "Missing encrypted_data field", http.StatusBadRequest)
			return
		}

		// Convert hexadecimal string to bytes
		cipherText, err := hex.DecodeString(encryptedData)
		if err != nil {
			http.Error(w, "Failed to decode hex", http.StatusInternalServerError)
			return
		}

		// Check if "aes_key" is present and the payload is larger than 512 bytes
		aesKeyHex, aesKeyPresent := encryptedJSON["aes_key"]
		if aesKeyPresent && len(cipherText) > 512 {
			// Decode the AES key from hexadecimal
			aesKeyBytes, err := hex.DecodeString(aesKeyHex)
			if err != nil {
				http.Error(w, "Failed to decode AES key hex", http.StatusInternalServerError)
				return
			}

			// Decrypt the AES key using the private key
			decryptedAESKey, err := rsakeys.DecryptWithPrivateKey(privateKey, aesKeyBytes)
			if err != nil {
				http.Error(w, "Failed to decrypt AES key", http.StatusInternalServerError)
				return
			}

			// Decrypt the encrypted data using the decrypted AES key
			decryptedData, err := rsakeys.DecryptAES(decryptedAESKey, cipherText)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Pass the decrypted JSON in the request body
			r.Body = io.NopCloser(bytes.NewReader([]byte(decryptedData)))
			r.ContentLength = int64(len(decryptedData))
			r.Header.Set("Content-Length", strconv.Itoa(len(decryptedData)))
		} else {
			// Decrypt the encrypted data using the private key
			decryptedData, err := rsakeys.DecryptWithPrivateKey(privateKey, cipherText)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Pass the decrypted JSON in the request body
			r.Body = io.NopCloser(bytes.NewReader([]byte(decryptedData)))
			r.ContentLength = int64(len(decryptedData))
			r.Header.Set("Content-Length", strconv.Itoa(len(decryptedData)))
		}

		next.ServeHTTP(w, r)
	})
}

// Helper function to extract the JWT token from the Authorization header
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Expecting the Authorization header value to be in the format "Bearer <token>"
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) == 2 && splitToken[0] == "Bearer" {
			return splitToken[1]
		}
	}
	return ""
}
