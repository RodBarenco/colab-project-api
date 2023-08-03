package main

import (
	"log"
	"os"

	"github.com/RodBarenco/colab-project-api/connection"
	"github.com/RodBarenco/colab-project-api/handlers"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Missing environment file path argument")
	}

	arg := os.Args[1]

	if err := handlers.InitHandlers(arg); err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	if err := handlers.JwtSecret(arg); err != nil {
		log.Fatalf("Failed to initialize JwtSecret: %v", err)
	}

	switch arg {
	case "1":
		connection.StartServer()

	case "2":
		connection.StartTestServer()
	default:
		log.Fatal("Was not possible to start Dev/Test mod!!!")
	}

}
