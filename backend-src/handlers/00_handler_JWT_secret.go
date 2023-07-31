package handlers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var jwtSecret = ""

func JwtSecret(envType string) error {
	var envFileName string
	switch envType {
	case "1":
		envFileName = ".env"
	case "2":
		envFileName = ".test.env"
	default:
		log.Fatalf("invalid envType: must be 1 or 2")
	}

	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatalf("failed to load %s file: %v", envFileName, err)
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET is not found in the environment")
	}

	jwtSecret = secret
	log.Println("Got JwtSecret!")
	return nil
}
