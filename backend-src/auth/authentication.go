package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
)

type LoginParams struct {
	Email    string
	Password string
}

func Login(DB *gorm.DB, params LoginParams, secret string) (string, int, error) {

	if params.Email == "" || params.Password == "" {
		return "", http.StatusBadRequest, fmt.Errorf("all required fields must be provided")
	}

	//lookup for the user
	user := db.User{}
	result := DB.First(&user, "email = ?", params.Email)

	if result.Error != nil {
		return "", http.StatusBadRequest, fmt.Errorf("User not found -  invalid email or password!")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))

	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("invalid email or password!")
	}

	//JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
		"aud": "user",
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("fail to creat token correctly")
	}

	return tokenString, http.StatusOK, nil
}
