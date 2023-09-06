package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
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

func ApproveArticleHandler(w http.ResponseWriter, r *http.Request) {
	// verify header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}
	var requestBody struct {
		ArticleID uuid.UUID `json:"ArticleID"`
		AdminID   uuid.UUID `json:"AdminID"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the user ID in the URL matches the one in the request body
	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.AdminID.String() {
		fmt.Print("url id: " + adminIDFromURL + " request body id: " + requestBody.AdminID.String())
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	// Parse the admin ID from the URL
	adminIDStr := chi.URLParam(r, "adminID")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid admin ID in URL")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	// get admin key
	adminPkey, err := db.GetAdminPublicKey(dbaccess, adminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// Call the ApproveArticle function
	message, err := db.ApproveArticle(dbaccess, requestBody.ArticleID, adminID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error approving article: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}
	message2 := "about article: " + requestBody.ArticleID.String() + ". - "

	// Respond with a success message
	RespondWithEncryptedJSON(w, http.StatusOK, message2+message, adminPkey)
}

func DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Verify header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		ArticleID uuid.UUID `json:"ArticleID"`
		AdminID   uuid.UUID `json:"AdminID"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the user ID in the URL matches the one in the request body
	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.AdminID.String() {
		fmt.Print("url id: " + adminIDFromURL + " request body id: " + requestBody.AdminID.String())
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	// Parse the admin ID from the URL
	adminIDStr := chi.URLParam(r, "adminID")
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid admin ID in URL")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	// get admin key
	adminPkey, err := db.GetAdminPublicKey(dbaccess, adminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// Call the DeleteArticle function
	err = db.DeleteArticle(dbaccess, requestBody.ArticleID, adminID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error deleting article: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := "about articcle: " + adminIDStr + " Article deleted successfully."
	// Respond with a success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, adminPkey)
}

func ApproveAdminHandler(w http.ResponseWriter, r *http.Request) {
	// Verify header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AdminIDToApprove uuid.UUID `json:"AdminIDToApprove"`
		ApprovingAdminID uuid.UUID `json:"ApprovingAdminID"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the user ID in the URL matches the one in the request body
	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.ApprovingAdminID.String() {
		fmt.Print("url id: " + adminIDFromURL + " request body id: " + requestBody.ApprovingAdminID.String())
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	// get admin key
	adminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.ApprovingAdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// Call the ApproveAdmin function
	err = db.ApproveAdmin(dbaccess, requestBody.AdminIDToApprove, requestBody.ApprovingAdminID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error approving admin: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := "admin with ID: " + requestBody.AdminIDToApprove.String() + " approved successfully."
	// Respond with a success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, adminPkey)
}

func DisapproveAdminHandler(w http.ResponseWriter, r *http.Request) {
	// Verify header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AdminIDToDisapprove uuid.UUID `json:"AdminIDToDisapprove"`
		DisapprovingAdminID uuid.UUID `json:"DisapprovingAdminID"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the user ID in the URL matches the one in the request body
	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.DisapprovingAdminID.String() {
		fmt.Print("url id: " + adminIDFromURL + " request body id: " + requestBody.DisapprovingAdminID.String())
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	// get admin key
	adminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.DisapprovingAdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	// Call the DisapproveAdmin function
	err = db.DisapproveAdmin(dbaccess, requestBody.AdminIDToDisapprove, requestBody.DisapprovingAdminID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error disapproving admin: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := "admin with ID: " + requestBody.AdminIDToDisapprove.String() + " disapproved successfully."
	// Respond with a success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, adminPkey)
}

func ModifyAdminPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	// Verify header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AllowerID   uuid.UUID `json:"AllowerID"`
		AllowedID   uuid.UUID `json:"AllowedID"`
		Permissions uint      `json:"Permissions"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	// Get the admin key of the allower
	allowerAdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.AllowerID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call the ModifyAdminPermissions function
	err = db.ModifyAdminPermissions(dbaccess, requestBody.AllowerID, requestBody.AllowedID, requestBody.Permissions)
	if err != nil {
		errorMessage := fmt.Sprintf("Error modifying admin permissions: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := "Admin with ID: " + requestBody.AllowedID.String() + " permissions modified successfully to: " + fmt.Sprint(requestBody.Permissions)
	// Respond with a success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, allowerAdminPkey)
}
