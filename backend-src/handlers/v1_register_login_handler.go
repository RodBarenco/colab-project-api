package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/RodBarenco/colab-project-api/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var body auth.SignupParams
	// Read the request body into the 'body' variable using json.NewDecoder
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create a slice to collect validation errors
	validationErrors := []string{}

	// Perform validation checks on the request data
	if !utils.IsValidFirstName(body.FirstName) {
		validationErrors = append(validationErrors, "First name must have 2 to 25 characters - and valid characters-words")
	}

	if !utils.IsValidLastName(body.LastName) {
		validationErrors = append(validationErrors, "Last name must have 1 to 40 characters - and valid characters-words")
	}

	if !utils.IsValidNickname(body.Nickname) {
		validationErrors = append(validationErrors, "First name must have 2 to 20 characters - and valid characters-words")
	}

	if !utils.IsValidEmail(body.Email) {
		validationErrors = append(validationErrors, "Invalid email format")
	}

	if !utils.IsValidPassword(body.Password) {
		validationErrors = append(validationErrors, "Password must have at least 5 characters - and valid characters-words")
	}

	// Check if the date of birth is valid (not greater than current date)
	if !utils.IsValidDateOfBirth(body.DateOfBirth) {
		validationErrors = append(validationErrors, "Invalid date of birth")
	}

	if !utils.IsValidField(body.Field) {
		validationErrors = append(validationErrors, "Field must have 2 to 50 characters - and valid characters-words")
	}

	if !utils.IsValidBiography(body.Biography) {
		validationErrors = append(validationErrors, "Biography must have 3 to 500 characters - and valid characters-words")
	}

	// Perform parallel database validations
	var wg sync.WaitGroup
	wg.Add(3)

	// Parallel validation for Ccourse
	go func() {
		defer wg.Done()
		if body.Ccourse != "" && body.CurrentlyID == nil {
			validationErrors = append(validationErrors, "CurrentlyID must be provided when Ccourse is provided")
			return
		}

		ccourseValidation := utils.IsValidCcourse(body.Ccourse, body.CurrentlyID, dbAccessor)
		if !ccourseValidation.IsValid {
			if !ccourseValidation.ExistsInDB {
				validationErrors = append(validationErrors, "Institution of current course does not exist or has not been registered")
			} else {
				validationErrors = append(validationErrors, "Invalid Ccourse format")
			}
		}
	}()

	// Parallel validation for Lcourse
	go func() {
		defer wg.Done()
		if body.Lcourse != "" && body.LastEducationID == nil {
			validationErrors = append(validationErrors, "LastEducationID must be provided when Lcourse is provided")
			return
		}

		lcourseValidation := utils.IsValidLcourse(body.Lcourse, body.LastEducationID, dbAccessor)
		if !lcourseValidation.IsValid {
			if !lcourseValidation.ExistsInDB {
				validationErrors = append(validationErrors, "Institution of last course does not exist or has not been registered")
			} else {
				validationErrors = append(validationErrors, "Invalid Lcourse format")
			}
		}
	}()

	// Parallel validation for interests
	go func() {
		defer wg.Done()
		isValidInterests, interestsValidationErr := utils.IsValidInterests(body.Interests, dbAccessor)
		if interestsValidationErr != nil {
			validationErrors = append(validationErrors, interestsValidationErr.Error())
		} else if !isValidInterests {
			validationErrors = append(validationErrors, "One or more interests are invalid or not registered")
		}
	}()

	wg.Wait()

	// CHECK VALIDATION ERRORS
	if len(validationErrors) > 0 {
		formattedErrors := make([]string, len(validationErrors))

		// Format each error message with the corresponding number
		for i, errMsg := range validationErrors {
			formattedErrors[i] = fmt.Sprintf("%d : %s", i+1, errMsg)
		}

		// Concatenate the formatted error messages into a single string
		errorMessage := strings.Join(formattedErrors, " , ")

		// Respond with the error message
		RespondWithError(w, http.StatusBadRequest, "{"+errorMessage+"}")
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
