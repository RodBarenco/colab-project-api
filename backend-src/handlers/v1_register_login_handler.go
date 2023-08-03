package handlers

import (
	"context"
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
		validationErrors = append(validationErrors, "Nickname must have 2 to 20 characters - and valid characters-words")
	}

	if !utils.IsValidEmail(body.Email) {
		validationErrors = append(validationErrors, "Invalid email format")
	}

	if !utils.IsValidPassword(body.Password) {
		validationErrors = append(validationErrors, "Password must have at least 5 characters - and valid characters-words")
	}

	if !utils.IsValidDateOfBirth(body.DateOfBirth) {
		validationErrors = append(validationErrors, "Invalid date of birth")
	}

	if !utils.IsValidField(body.Field) {
		validationErrors = append(validationErrors, "Field must have 2 to 50 characters - and valid characters-words")
	}

	if !utils.IsValidBiography(body.Biography) {
		validationErrors = append(validationErrors, "Biography must have 3 to 500 characters - and valid characters-words")
	}

	// Perform parallel database validations - that need to access DB
	var wg sync.WaitGroup
	wg.Add(3)

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

		errorMessage := strings.Join(formattedErrors, " , ")

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
		RespondWithError(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}

	// Call the Signup function passing the validated params and the database connection
	statusCode, err := auth.Signup(db, params)
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
		},
		Message: "User registered!",
	}

	RespondWithJSON(w, http.StatusCreated, response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	db := dbAccessor

	// Call the Login function passing the validated params and the database connection
	tokenString, statusCode, err := auth.Login(db, body, jwtSecret)
	if err != nil {
		errorMessage := fmt.Sprintf("Error during login: %v", err)
		RespondWithError(w, statusCode, errorMessage)
		return
	}

	// Set the token in the request context
	ctx := context.WithValue(r.Context(), "jwtToken", tokenString)

	// Call the next handler with the updated context.
	r = r.WithContext(ctx)

	response := res.LoginRes{
		Message: "Login successful!",
		Token:   tokenString,
	}

	RespondWithJSON(w, statusCode, response)
}
