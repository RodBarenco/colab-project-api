package handlers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/joho/godotenv"
)

var dbAccessor *gorm.DB

// InitDB initializes the database connection for the handlers package
func InitHandlers(envType string) error {
	var envFileName string
	switch envType {
	case "1":
		envFileName = ".env"
	case "2":
		envFileName = ".test.env"
	default:
		return fmt.Errorf("invalid envType: must be 1 or 2")
	}

	err := godotenv.Load(envFileName)
	if err != nil {
		log.Fatalf("failed to load %s file: %v", envFileName, err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	dbAccessInstance, err := db.DBaccess(dsn)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	dbAccessor = dbAccessInstance.DB
	log.Println("Handlers connected successfully!")
	return nil
}
