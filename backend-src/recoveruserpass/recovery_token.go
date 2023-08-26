package recoveruserpass

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RodBarenco/colab-project-api/rsakeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secretOfTheKey = ""

func InitRecoverTokenSecret(arg string) error {
	var envFileName string
	switch arg {
	case "1":
		envFileName = ".env"
	case "2":
		envFileName = ".test.env"
	default:
		return fmt.Errorf("invalid recover token env type: must be 1 or 2")
	}
	err := godotenv.Load(envFileName)
	if err != nil {
		log.Printf("failed to load %s file: %v", envFileName, err)
	}

	secretOfTheKey = os.Getenv("SECRET")
	if secretOfTheKey == "" {
		log.Fatal("SECRET was not found in the environment")
	}
	return nil
}

// secret key
var jwtSecret = []byte(secretOfTheKey)

// GenerateRecoveryToken gera um token de recuperação de senha com validade de 10 minutos.
func GenerateRecoveryToken(randomPass string) (string, error) {
	// Obter a chave pública usada para criptografia
	publicKey, err := rsakeys.ReadPublicKeyFromFile("public_key.der")
	if err != nil {
		return "", err
	}

	// Criptografar a senha aleatória
	encryptedPass, err := rsakeys.EncryptAndEncode(publicKey, []byte(randomPass))
	if err != nil {
		return "", err
	}

	// Definir os claims do token
	claims := jwt.MapClaims{
		"exp":        time.Now().Add(time.Minute * 10).Unix(), // Token expira em 10 minutos
		"randomPass": encryptedPass,                           // Incluir a senha aleatória criptografada no token
	}

	// Criar token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assinar o token com a chave secreta
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
