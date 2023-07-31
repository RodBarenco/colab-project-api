package connection

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/routes"
	"github.com/joho/godotenv"
)

func StartTestServer() {
	err := godotenv.Load(".test.env")
	if err != nil {
		panic(fmt.Errorf("failed to load .test.env file: %w", err))
	}

	// Connecting with PostgreSQL
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbAccessInstance, err := db.DBaccess(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Errorf("failed to connect to the database: %w", err))
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
	log.Printf("\033[33mTestServer starting on PORT: %v\033[0m", portString)
	log.Println("Run your tests against this server.")
	log.Println("Remember to properly clean up the resources when testing is done.")

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
