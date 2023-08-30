package handlers

import (
	"encoding/json"
	"log"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/utils"
	"gorm.io/gorm"
)

func GenerateRootAdminIfNeededHandler() {
	// Check if there is at least one admin in the database
	accessor := dbAccessor

	var admin db.Admin
	result := accessor.First(&admin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			params := auth.AdminRegistrationParams{
				FirstName:   "Root",
				LastName:    "Admin",
				Nickname:    "001",
				Email:       "root@admin.com",
				Password:    "yourpassword",
				DateOfBirth: "2000-01-01",
				Permissions: 4,
				IsAccepted:  true,
			}

			message, pbKeyStr, err := auth.RegisterAdmin(params, accessor)
			if err != nil {
				log.Panic(utils.RedColor.InitColor+"Failed to generate root admin: %v"+utils.EndColor, err)
			}
			params.PublicKey = pbKeyStr

			jsonData, err := json.MarshalIndent(params, "", "  ")
			if err != nil {
				log.Fatalf("Erro ao converter para JSON: %v", err)
			}

			log.Printf(utils.GreenColor.InitColor+"\nRoot admin generated successfully: %v"+utils.EndColor, message)
			log.Printf(utils.OrangeColor.InitColor + "\n Remamber to change the fowlling fields by login in your Admin account with your given email and password: " + utils.EndColor)
			log.Printf(string(jsonData))

		} else {
			log.Panic(utils.RedColor.InitColor+"Failed to check admin existence: %v"+utils.EndColor, result.Error)
		}
	}
	log.Printf("Everything is right to start the server...")
	return
}
