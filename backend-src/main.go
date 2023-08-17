package main

import (
	"log"
	"os"

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

	StartApp(arg)
}
