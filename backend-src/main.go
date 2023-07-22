package main

import (
	"log"
	"os"

	"github.com/RodBarenco/colab-project-api/connection"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Missing environment file path argument")
	}

	arg := os.Args[1]

	switch arg {
	case "1":
		connection.StartServer()

	case "2":
		connection.StartTest()

	default:
		log.Fatal("Was not possible to start Dev/Test mod!!!")
	}

}
