package auth

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
)

// SignupParams holds the parameters required for user registration.
type SignupParams struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	DateOfBirth time.Time
}

// Signup creates a new user and saves it to the database.
func Signup(ctx context.Context, DB *gorm.DB, params SignupParams) error {
	// Check if required fields are not empty
	if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.Password == "" || params.DateOfBirth.IsZero() {
		return fmt.Errorf("all required fields must be provided")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		// Return a more specific error message including the original password
		return fmt.Errorf("failed to hash the password: %w", err)
	}

	// Create the user with the hashed password
	user := db.User{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		Password:    string(hashedPassword),
		DateOfBirth: params.DateOfBirth,
	}

	// Save the user to the database using gorm's Create method
	result := DB.Create(&user)
	if result.Error != nil {
		return fmt.Errorf("failed to create the user: %w", result.Error)
	}

	return nil
}
