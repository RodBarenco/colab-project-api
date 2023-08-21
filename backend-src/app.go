package main

import (
	"log"

	"time"

	"github.com/RodBarenco/colab-project-api/connection"
	"github.com/RodBarenco/colab-project-api/rsakeys"
)

func StartApp(arg string) {
	validUntil := time.Now().Add(7 * 24 * time.Hour)

	if err := rsakeys.EnsureKeysValid(validUntil); err != nil {
		log.Fatalf("Failed to ensure RSA keys are valid: %v", err)
	}

	switch arg {
	case "1":
		connection.StartServer()
	case "2":
		connection.StartTestServer()
	default:
		log.Fatal("Was not possible to start Dev/Test mode!!!")
	}
}
