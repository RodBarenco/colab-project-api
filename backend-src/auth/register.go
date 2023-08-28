package auth

import (
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/google/uuid"
)

// SignupParams holds the parameters required for user registration.
type SignupParams struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	DateOfBirth     string
	Nickname        string
	Field           string
	Interests       []*db.Interest
	Biography       string
	Lcourse         string
	Ccourse         string
	LastEducationID *uuid.UUID
	CurrentlyID     *uuid.UUID
	OpenToColab     bool
	CreatedAt       time.Time
	ProfilePhoto    string
}

// Signup creates a new user and saves it to the database.
func Signup(DB *gorm.DB, params SignupParams) (int, error) {
	// Check if required fields are not empty
	if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.Password == "" || params.DateOfBirth == "" {
		return http.StatusBadRequest, fmt.Errorf("all required fields must be provided")
	}

	// Parse the date of birth to the desired format (YYYY-MM-DD)
	dateOfBirth, err := time.Parse("2006-01-02", params.DateOfBirth)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid date of birth format. It must be in the format YYYY-MM-DD")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		// Return a more specific error message including the original password
		return http.StatusInternalServerError, fmt.Errorf("failed to hash the password: %w", err)
	}

	// Create the user with the hashed password
	user := db.User{
		FirstName:       params.FirstName,
		LastName:        params.LastName,
		Email:           params.Email,
		Password:        hashedPassword,
		DateOfBirth:     dateOfBirth,
		Nickname:        params.Nickname,
		Field:           params.Field,
		Interests:       params.Interests,
		Biography:       params.Biography,
		Lcourse:         params.Lcourse,
		Ccourse:         params.Ccourse,
		LastEducationID: params.LastEducationID,
		CurrentlyID:     params.CurrentlyID,
		OpenToColab:     params.OpenToColab,
		ProfilePhoto:    params.ProfilePhoto,
	}

	// Save the user to the database using gorm's Create method
	result := DB.Create(&user)
	if result.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create the user: %w", result.Error)
	}

	return http.StatusCreated, nil
}
