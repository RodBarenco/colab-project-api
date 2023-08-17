package main

import (
	"log"

	"time"

	"github.com/RodBarenco/colab-project-api/connection"
	"github.com/RodBarenco/colab-project-api/rsakeys"
)

func StartApp(arg string) {
	// heck intervalc
	checkInterval := 24 * time.Hour

	for {
		validUntil := time.Now().Add(7 * 24 * time.Hour)

		if err := rsakeys.EnsureKeysValid(validUntil); err != nil {
			log.Fatalf("Failed to ensure RSA keys are valid: %v", err)
		}

		// Inicia o servidor conforme o argumento passado
		switch arg {
		case "1":
			connection.StartServer()

		case "2":
			connection.StartTestServer()
		default:
			log.Fatal("Was not possible to start Dev/Test mod!!!")
		}

		// Aguarda o próximo intervalo antes de verificar novamente
		time.Sleep(checkInterval)
	}
}
