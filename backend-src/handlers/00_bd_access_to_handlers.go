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
func InitHandlers() error {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env:", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	dbAccessInstance, err := db.DBaccess(dsn)
	if err != nil {
		return err
	}

	dbAccessor = dbAccessInstance.DB
	fmt.Println("Handlers connected successfully!")
	return nil
}
