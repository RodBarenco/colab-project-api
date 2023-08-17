package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/RodBarenco/colab-project-api/rsakeys"
)

type LoginParams struct {
	Email    string
	Password string
}

func Login(DB *gorm.DB, params LoginParams, secret string) (res.LoginRes, int, error) {
	var loginRes res.LoginRes

	if params.Email == "" || params.Password == "" {
		return loginRes, http.StatusBadRequest, fmt.Errorf("all required fields must be provided")
	}

	// Lookup for the user
	user := db.User{}
	result := DB.First(&user, "email = ?", params.Email)

	if result.Error != nil {
		return loginRes, http.StatusBadRequest, fmt.Errorf("User not found - invalid email or password!")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return loginRes, http.StatusBadRequest, fmt.Errorf("invalid email or password!")
	}

	// JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
		"aud": "user",
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return loginRes, http.StatusBadRequest, fmt.Errorf("failed to create token correctly")
	}

	// Read the public key from the file
	publicKey, err := rsakeys.ReadPublicKeyFromFile("public_key.der")
	if err != nil {
		return loginRes, http.StatusInternalServerError, fmt.Errorf("failed to read public key: %v", err)
	}

	// Populate the LoginRes struct with the necessary data
	loginRes.Message = "Login successful!"
	loginRes.Token = tokenString
	loginRes.PublicKey = publicKey
	loginRes.UserID = user.ID

	return loginRes, http.StatusOK, nil
}
