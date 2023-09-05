package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	LoginRes, code, err := GenericLogin(DB, "user", params, secret)
	if err != nil {
		return LoginRes, code, err
	}
	return LoginRes, code, nil
}

func AdminLogin(DB *gorm.DB, params LoginParams, secret string) (res.LoginRes, int, error) {
	LoginRes, code, err := GenericLogin(DB, "admin", params, secret)
	if err != nil {
		return LoginRes, code, err
	}
	return LoginRes, code, nil
}

func GenericLogin(DB *gorm.DB, userType string, params LoginParams, secret string) (res.LoginRes, int, error) {
	var loginRes res.LoginRes

	if params.Email == "" || params.Password == "" {
		return loginRes, http.StatusBadRequest, fmt.Errorf("all required fields must be provided")
	}

	// variables
	user := db.User{}
	admin := db.Admin{}
	var pass string
	var yourID uuid.UUID
	var result *gorm.DB

	if userType == "user" {
		result = DB.First(&user, "email = ?", params.Email)
		pass = user.Password
		yourID = user.ID
	} else if userType == "admin" {
		result = DB.First(&admin, "email = ?", params.Email)
		pass = admin.Password
		yourID = admin.ID
	} else {
		return loginRes, http.StatusBadRequest, fmt.Errorf("invalid user type")
	}

	if result.Error != nil {
		return loginRes, http.StatusBadRequest, fmt.Errorf("%s not found - invalid email or password!", userType)
	}

	err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(params.Password))
	if err != nil {
		return loginRes, http.StatusBadRequest, fmt.Errorf("invalid email or password!")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": yourID, // or user.(*db.Admin).ID, based on userType
		"exp": time.Now().Add(time.Hour * 12).Unix(),
		"aud": userType,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return loginRes, http.StatusBadRequest, fmt.Errorf("failed to create token correctly")
	}

	publicKey, err := rsakeys.ReadPublicKeyFromFile("public_key.der")
	if err != nil {
		return loginRes, http.StatusInternalServerError, fmt.Errorf("failed to read public key: %v", err)
	}

	loginRes.Message = "Login successful!"
	loginRes.Token = tokenString
	loginRes.PublicKey = publicKey
	loginRes.YourID = yourID

	return loginRes, http.StatusOK, nil
}
