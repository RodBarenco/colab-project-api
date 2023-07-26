package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/res"
)

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var body auth.SignupParams
	// Read the request body into the 'body' variable using json.NewDecoder
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Perform validation checks on the request data
	if len(body.FirstName) < 2 {
		RespondWithError(w, http.StatusBadRequest, "First name must have at least 2 characters")
		return
	}

	if len(body.LastName) < 1 {
		RespondWithError(w, http.StatusBadRequest, "Last name must have at least 2 characters")
		return
	}

	if !isValidEmail(body.Email) {
		RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	if len(body.Password) < 5 {
		RespondWithError(w, http.StatusBadRequest, "Password must have at least 5 characters")
		return
	}

	// Check if the date of birth is valid (not greater than current date)
	if body.DateOfBirth.After(time.Now()) {
		RespondWithError(w, http.StatusBadRequest, "Invalid date of birth")
		return
	}

	// Validate that if Ccourse is provided, CurrentlyID must also be provided
	if body.Ccourse != "" && body.CurrentlyID == nil {
		RespondWithError(w, http.StatusBadRequest, "Ccourse can only be added if CurrentlyID is provided")
		return
	}

	// Validate that if Lcourse is provided, LastEducationID must also be provided
	if body.Lcourse != "" && body.LastEducationID == nil {
		RespondWithError(w, http.StatusBadRequest, "Lcourse can only be added if LastEducationID is provided")
		return
	}

	// Convert the 'body' data into SignupParams
	params := auth.SignupParams{
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		Email:           body.Email,
		Password:        body.Password,
		DateOfBirth:     body.DateOfBirth,
		Nickname:        body.Nickname,
		Field:           body.Field,
		Interests:       body.Interests,
		Biography:       body.Biography,
		Lcourse:         body.Lcourse,
		Ccourse:         body.Ccourse,
		LastEducationID: body.LastEducationID,
		CurrentlyID:     body.CurrentlyID,
		OpenToColab:     body.OpenToColab,
	}

	// Access the gorm.DB connection using dbAccessor
	db := dbAccessor

	if err != nil {
		// Handle the error if there's an issue with the database connection
		RespondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}

	// Call the Signup function passing the validated params and the database connection
	statusCode, err := auth.Signup(r.Context(), db, params)
	if err != nil {

		errorMessage := fmt.Sprintf("Error during signup: %v", err)
		RespondWithError(w, statusCode, errorMessage)
		return
	}

	// If signup is successful, respond with a success JSON
	response := res.SignupRes{
		User: res.UserGetedResponse{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Email:     body.Email,
			// Add other fields as needed
		},
		Message: "User registered!",
	}

	// Respond with the SignupRes JSON
	RespondWithJSON(w, http.StatusCreated, response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the logic for user login.
}
