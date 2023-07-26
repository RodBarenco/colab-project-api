package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/google/uuid"
)

// SignupParams holds the parameters required for user registration.
type SignupParams struct {
	FirstName       string
	LastName        string
	Email           string
	Password        string
	DateOfBirth     time.Time
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
}

// Signup creates a new user and saves it to the database.
func Signup(ctx context.Context, DB *gorm.DB, params SignupParams) (int, error) {
	// Check if required fields are not empty
	if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.Password == "" || params.DateOfBirth.IsZero() {
		return http.StatusBadRequest, fmt.Errorf("all required fields must be provided")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		// Return a more specific error message including the original password
		return http.StatusInternalServerError, fmt.Errorf("failed to hash the password: %w", err)
	}

	// Create the user with the hashed password
	user := db.User{
		FirstName:       params.FirstName,
		LastName:        params.LastName,
		Email:           params.Email,
		Password:        string(hashedPassword),
		DateOfBirth:     params.DateOfBirth,
		Nickname:        params.Nickname,
		Field:           params.Field,
		Interests:       params.Interests,
		Biography:       params.Biography,
		Lcourse:         params.Lcourse,
		Ccourse:         params.Ccourse,
		LastEducationID: params.LastEducationID,
		CurrentlyID:     params.CurrentlyID,
		OpenToColab:     params.OpenToColab,
	}

	// Save the user to the database using gorm's Create method
	result := DB.Create(&user)
	if result.Error != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create the user: %w", result.Error)
	}

	return http.StatusCreated, nil
}
