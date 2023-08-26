package recoveruserpass

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func GenerateToken(randomPass string) (string, error) {

	tokenString, err := GenerateRecoveryToken(randomPass)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SendRecoveryEmail(email string) (string, error) {
	serv, mail, key, errSmtp := LoadSMTPInfo()
	if errSmtp != nil {
		return "", errSmtp
	}
	randomPass := generateRandomPassword(8)
	token, errGt := GenerateToken(randomPass)
	if errGt != nil {
		return "cannot generate token", errGt
	}
	// Configurar a mensagem de e-mail
	msg := gomail.NewMessage()
	msg.SetHeader("From", "projectcolab.br@gamil.com")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Recuperação de Senha - Projeto Colab")
	msg.SetBody("text/plain", "Olá! \n Aqui está o link para recuperação de senha: \n http://localhost:8000/v1/recoverynow/?token="+token+"\n"+"Sua senha de acesso é "+randomPass) // modify this when in produciton

	// Configurar o cliente SMTP e enviar o e-mail
	mailer := gomail.NewDialer(serv, 1025, mail, key) // you may need to chagenge the port, problably to 587
	err := mailer.DialAndSend(msg)
	if err != nil {
		return "couldn't sent recovery email.", err
	}
	return "recovery email sent successfully.", nil
}

type RecoverParams struct {
	UserID     uuid.UUID
	Email      string
	Password   string
	RandomPass string
}

func UpdatePassword(accessor *gorm.DB, userID uuid.UUID, email, newPassword, token string) error {

	user := db.User{ID: userID}

	if newPassword != "" && !utils.IsValidPassword(newPassword) {
		return fmt.Errorf("invalid password!")
	}
	user.Password = newPassword

	return nil
}

// HELPER ------------------------------------------------------------

func ValidateRecoveryToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid or expired JWT token")
	}

	// Verificar randomPass
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to get the token claims")
	}

	randomPassHex, found := claims["randomPass"].(string)
	if !found {
		return "", fmt.Errorf("randomPass not found in token")
	}

	// Agora você pode utilizar o valor de randomPassHex como necessário

	return randomPassHex, nil
}

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomBytes)
}

// LoadSMTPInfo loads SMTP configuration from the .env file.
// It returns the SMTP server, sender email, API key, and any error encountered.
func LoadSMTPInfo() (string, string, string, error) {
	// Load variables from the .env file
	if err := godotenv.Load(".env"); err != nil {
		return "", "", "", err
	}

	apiKey := os.Getenv("SMTP_API_KEY")
	smtpServer := os.Getenv("SMTP_SERVER")
	senderEmail := os.Getenv("SMTP_SENDER")

	return smtpServer, senderEmail, apiKey, nil
}
