package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminRegistrationParams struct {
	FirstName    string
	LastName     string
	Nickname     string
	Email        string
	Password     string
	DateOfBirth  string
	Title        string
	Field        string
	Biography    string
	CurrentlyID  *uuid.UUID
	PublicKey    string
	IsAccepted   bool
	Permissions  uint
	ProfilePhoto string
}

// this function recives params for registration, and the access to database and generat a Admin
// then returns []friendly messages, a publicKey string and error
func RegisterAdmin(params AdminRegistrationParams, accessor *gorm.DB) ([]string, string, error) {
	cB := utils.BlueColor.InitColor
	cY := utils.YellowColor.InitColor
	cC := utils.CyanColor.InitColor
	rst := utils.EndColor

	messages := []string{}
	// Check if required fields are not empty
	if params.FirstName == "" || params.LastName == "" || params.Email == "" || params.Password == "" ||
		params.DateOfBirth == "" || params.Permissions < 0 || params.Permissions > 4 {
		return messages, "", fmt.Errorf("all required fields must be provided")
	}
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		messages = []string{fmt.Sprintf("error generating RSA key pair: %v", err)}
		return messages, "", err
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKeyBytes)
	messages = append(messages, fmt.Sprintf("\n %v ATENTION:%v %v  This is your private key. Under no circumstances lose this key. Take all security measures to protect it, that's your responsibility. Only with it you will be able to decode the messages you receive: %v \n %v", cY, rst, cB, rst, privateKeyBase64))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		messages = []string{fmt.Sprintf("error exporting public key: %v", err)}
		return messages, "", err
	}

	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)
	messages = append(messages, fmt.Sprintf("\n%v This is your public key and it will be exposed every time you log in: %v \n %v", cC, rst, publicKeyBase64))

	if err != nil {
		messages = []string{fmt.Sprintf("failed to generate publicKeyBase64 response")}
		return messages, "", err
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		messages = []string{fmt.Sprintf("error hashing password: %v", err)}
		return messages, "", err
	}

	dateOfBirth, err := time.Parse("2006-01-02", params.DateOfBirth)
	if err != nil {
		messages = []string{fmt.Sprintf("invalid date of birth format. It must be in the format YYYY-MM-DD")}
		return messages, "", err
	}

	// Create the admin record
	admin := db.Admin{
		ID:           uuid.New(),
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Nickname:     params.Nickname,
		Email:        params.Email,
		Password:     string(hashedPassword),
		DateOfBirth:  dateOfBirth,
		Title:        params.Title,
		Field:        params.Field,
		Biography:    params.Biography,
		CurrentlyID:  params.CurrentlyID,
		CreatedAt:    time.Now(),
		PublicKey:    publicKeyBase64,
		Permissions:  params.Permissions,
		ProfilePhoto: params.ProfilePhoto,
		IsAccepted:   params.IsAccepted,
	}

	// Save the admin record to the database using GORM
	result := accessor.Create(&admin)
	if result.Error != nil {
		messages = []string{fmt.Sprintf("failed to create Admin: %v", result.Error)}
		return messages, "", err
	}

	// Return a success message
	return messages, publicKeyBase64, nil
}
