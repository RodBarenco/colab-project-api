package connection

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/routes"
	"github.com/joho/godotenv"
)

func StartServer() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env:", err)
	}

	// Connecting with PostgreSQL
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not found in the environment")
	}

	// Use DBaccess function from the db package to get the dbAccess instance
	dbAccessInstance, err := db.DBaccess(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Access the gorm.DB connection from the dbAccess instance
	GormDB := dbAccessInstance.DB

	// Call migration from the db package, passing the gorm.DB connection
	err = db.Migrate(GormDB)
	if err != nil {
		log.Fatal("Failed to perform migration:", err)
	}

	// Call main router from the routes package
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET is not found in the environment")
	}

	router := routes.MainRouter(secret)

	// Start the server...
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
