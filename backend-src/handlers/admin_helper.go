package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/rsakeys"
	"github.com/RodBarenco/colab-project-api/utils"
)

func GetPKeyHandler(w http.ResponseWriter, r *http.Request) {
	publicKey, err := rsakeys.ReadPublicKeyFromFile("public_key.der")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := "Use this publick key: " + publicKey

	RespondWithJSON(w, http.StatusOK, response)
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	// verify header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var body auth.LoginParams
	// Read the request body into the 'body' variable using json.NewDecoder
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validations
	if !utils.IsValidEmail(body.Email) {
		RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	if !utils.IsValidPassword(body.Password) {
		RespondWithError(w, http.StatusBadRequest, "Password must have at least 5 characters - and valid characters-words")
		return
	}

	// Access the gorm.DB connection using dbAccessor
	access := dbAccessor

	// Call the Login function passing the validated params and the database connection
	resLog, statusCode, err := auth.AdminLogin(access, body, jwtSecret)
	if err != nil {
		errorMessage := fmt.Sprintf("Error during login: %v", err)
		RespondWithError(w, statusCode, errorMessage)
		return
	}

	AdminID, err := utils.GetAdminIDByEmail(access, body.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	adminPkey, err := db.GetAdminPublicKey(access, AdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// Set the token in the request context
	ctx := context.WithValue(r.Context(), "jwtToken", resLog.Token)

	// Call the next handler with the updated context.
	r = r.WithContext(ctx)

	response := resLog

	RespondWithEncryptedJSON(w, statusCode, response, adminPkey)
}
