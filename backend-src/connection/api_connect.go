package connection

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/routes"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartServer() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env:", err)
	}

	// connecting with PostgreSQL
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Call migration
	err = db.Migrate(gormDB)
	if err != nil {
		log.Fatal("Failed to perform migration:", err)
	}

	//call main router
	router := routes.MainRouter()

	// start server...
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Log the server starting message
	log.Printf("\033[33mServer starting on PORT: %v\033[0m", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
