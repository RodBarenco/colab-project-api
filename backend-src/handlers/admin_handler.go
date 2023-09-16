package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
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

func RegisterNewAdmin(w http.ResponseWriter, r *http.Request) {
	var body auth.NewAdminRegistrationParams

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	//validations ---------------------------------------------------

	if !utils.IsValidFirstName(body.FirstName) {
		RespondWithError(w, http.StatusBadRequest, "First name must have 2 to 25 characters - and valid characters-words")
		return
	}

	if !utils.IsValidLastName(body.LastName) {
		RespondWithError(w, http.StatusBadRequest, "Last name must have 1 to 40 characters - and valid characters-words")
		return
	}

	if !utils.IsValidNickname(body.Nickname) {
		RespondWithError(w, http.StatusBadRequest, "Nickname must have 2 to 30 characters - and valid characters-words")
		return
	}

	if !utils.IsValidEmail(body.Email) {
		RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	if !utils.IsValidPassword(body.Password) {
		RespondWithError(w, http.StatusBadRequest, "Password must have at least 5 characters - and valid characters-words")
		return
	}

	if !utils.IsValidDateOfBirth(body.DateOfBirth) {
		RespondWithError(w, http.StatusBadRequest, "Invalid Data of Birth")
		return
	} else {
		dob, _ := time.Parse("2006-01-02", body.DateOfBirth)

		now := time.Now()

		if dob.After(now) {
			RespondWithError(w, http.StatusBadRequest, "Date of birth cannot be in the future")
			return
		}
	}

	if !utils.IsValidField(body.Field) {
		RespondWithError(w, http.StatusBadRequest, "Field must have 2 to 50 characters - and valid characters-words")
		return
	}

	if !utils.IsValidBiography(body.Biography) {
		RespondWithError(w, http.StatusBadRequest, "Biography must have 3 to 500 characters - and valid characters-words")
		return
	}

	if !utils.IsValidTitle(body.Title) {
		RespondWithError(w, http.StatusBadRequest, "Title must have 3 to 50 characters - and valid characters-words")
		return
	}

	if !utils.IsValidImage(body.ProfilePhoto) {
		RespondWithError(w, http.StatusBadRequest, "Invalid photo format")
		return
	}

	imageURL, err := SaveImageToDBHandler(body.ProfilePhoto) // here saves the photo
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if there is at least one admin in the database
	accessor := dbAccessor

	params := auth.AdminRegistrationParams{
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		Nickname:     body.Nickname,
		Email:        body.Email,
		Password:     body.Password,
		DateOfBirth:  body.DateOfBirth,
		Field:        body.Field,
		Biography:    body.Biography,
		ProfilePhoto: imageURL,
		Title:        body.Title,
		Permissions:  0,
		IsAccepted:   false,
	}

	message1, pbKeyStr, err := auth.RegisterAdmin(params, accessor)
	if err != nil {
		log.Panic(utils.RedColor.InitColor+"Failed to generate root admin: %v"+utils.EndColor, err)
	}
	params.PublicKey = pbKeyStr

	jsonData, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		log.Fatalf("Erro ao converter para JSON: %v", err)
	}

	log.Printf(utils.GreenColor.InitColor+"\n Admin registred successfully: %v"+utils.EndColor, message1)
	log.Printf(string(jsonData))

	message2 := res.SignupRes{
		User: res.UserGetedResponse{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Email:     body.Email,
		},
		Message: "Admin fields sent, if everything is ok with the first message you should just wait for approval. Again, remember to save your private key, and in a secure location, it's unique and we don't have access to it.",
	}

	response := res.AdminSignup{
		Message1: message1,
		Message2: message2,
	}

	RespondWithJSON(w, http.StatusCreated, response)
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

	// Check if the user ID in the URL matches the one in the request body
	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.AllowerID.String() {
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
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

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		UserIDToDelete  uuid.UUID `json:"UserIDToDelete"`
		DeletingAdminID uuid.UUID `json:"DeletingAdminID"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.DeletingAdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	dbaccess := dbAccessor

	deletingAdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.DeletingAdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// call DeleteUser
	err = db.DeleteUser(dbaccess, requestBody.UserIDToDelete, requestBody.DeletingAdminID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error deleting user: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := "User with ID: " + requestBody.UserIDToDelete.String() + " deleted successfully."

	RespondWithEncryptedJSON(w, http.StatusOK, response, deletingAdminPkey)
}

func DeleteAdminHandler(w http.ResponseWriter, r *http.Request) {

	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AdminIDToDelete uuid.UUID `json:"AdminIDToDelete"`
		DeletingAdminID uuid.UUID `json:"DeletingAdminID"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.DeletingAdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	dbaccess := dbAccessor

	deletingAdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.DeletingAdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// call DeleteAdmin
	err = db.DeleteAdmin(dbaccess, requestBody.AdminIDToDelete, requestBody.DeletingAdminID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error deleting admin: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := "Admin with ID: " + requestBody.AdminIDToDelete.String() + " deleted successfully."

	RespondWithEncryptedJSON(w, http.StatusOK, response, deletingAdminPkey)
}

func GetUnacceptedArticlesIDHandler(w http.ResponseWriter, r *http.Request) {
	// Check the header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AdminID       uuid.UUID `json:"AdminID"`
		DateToCompare string    `json:"DateToCompare"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the root admin ID in the URL matches the ID in the request body
	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.AdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	// parse string data "2006-01-02" to time.Time
	dateToCompare, err := time.Parse("2006-01-02", requestBody.DateToCompare)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid date format. Please use the format '2006-01-02'")
		return
	}

	dbaccess := dbAccessor

	AdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// call GetUnacceptedArticlesID
	unacceptedArticleIDs, err := db.GetUnacceptedArticlesID(dbaccess, dateToCompare)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting unaccepted article IDs: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	articles, err := db.GetUnacceptedArticlesByID(dbaccess, unacceptedArticleIDs)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error geting unaccepted articles by id")
		return
	}

	var unacceptedArticles []res.ArticleAdminFormat

	for _, article := range articles {
		articleFormat := res.ArticleAdminFormat{
			ID:             article.ID,
			Title:          article.Title,
			Field:          article.Field,
			SubmissionDate: article.SubmissionDate.Format("2006-01-02"),
		}
		unacceptedArticles = append(unacceptedArticles, articleFormat)
	}

	response := res.ArticleAdminFormatRes{
		Articles: unacceptedArticles,
		Message:  "Unaccepted articles retrieved successfully",
	}

	RespondWithEncryptedJSON(w, http.StatusOK, response, AdminPkey)
}

func GetUnacceptedArticlesByFieldHandler(w http.ResponseWriter, r *http.Request) {

	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AdminID        uuid.UUID `json:"AdminID"`
		FieldToCompare string    `json:"FieldToCompare"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	adminIDFromURL := chi.URLParam(r, "adminID")
	if adminIDFromURL != requestBody.AdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Admin IDs do not match")
		return
	}

	dbaccess := dbAccessor

	AdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Call GetUnacceptedArticlesIDByField
	unacceptedArticleIDs, err := db.GetUnacceptedArticlesIDByField(dbaccess, requestBody.FieldToCompare)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting unaccepted article IDs by field: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	// call GetUnacceptedArticlesByID
	articles, err := db.GetUnacceptedArticlesByID(dbaccess, unacceptedArticleIDs)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error geting unaccepted articles by id")
		return
	}

	var unacceptedArticles []res.ArticleAdminFormat

	for _, article := range articles {
		articleFormat := res.ArticleAdminFormat{
			ID:             article.ID,
			Title:          article.Title,
			Field:          article.Field,
			SubmissionDate: article.SubmissionDate.Format("2006-01-02"),
		}
		unacceptedArticles = append(unacceptedArticles, articleFormat)
	}

	// build response
	response := res.ArticleAdminFormatRes{
		Articles: unacceptedArticles,
		Message:  "Unaccepted articles retrieved successfully",
	}

	RespondWithEncryptedJSON(w, http.StatusOK, response, AdminPkey)
}

func CleanOldUnacceptedArticlesByDateHandler(w http.ResponseWriter, r *http.Request) {
	// Check the header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	// Decode the request body
	var requestBody struct {
		AdminID       uuid.UUID `json:"AdminID"`
		DateToCompare string    `json:"DateToCompare"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the root admin ID in the URL matches the ID in the request body
	rootAdminIDFromURL := chi.URLParam(r, "rootAdminID")
	if rootAdminIDFromURL != requestBody.AdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Root admin IDs do not match")
		return
	}

	// string to date convertion
	dateToCompare, err := time.Parse("2006-01-02", requestBody.DateToCompare)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid date format. Please use the format '2006-01-02'")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	AdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.CleanOldUnacceptedArticlesByDate(dbaccess, dateToCompare, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to delete articles: "+err.Error())
		return
	}

	// Build the response
	response := "old unaccepted articles cleaned successfully until " + requestBody.DateToCompare

	// Respond with the success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, AdminPkey)
}

func CleanOldUnacceptedArticlesByDateAndFieldHandler(w http.ResponseWriter, r *http.Request) {
	// Check the header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	// Decode the request body
	var requestBody struct {
		AdminID       uuid.UUID `json:"AdminID"`
		DateToCompare string    `json:"DateToCompare"`
		Field         string    `json:"Field"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the root admin ID in the URL matches the ID in the request body
	rootAdminIDFromURL := chi.URLParam(r, "rootAdminID")
	if rootAdminIDFromURL != requestBody.AdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Root admin IDs do not match")
		return
	}

	dateToCompare, err := time.Parse("2006-01-02", requestBody.DateToCompare)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid date format. Please use the format '2006-01-02'")
		return
	}

	//validate field
	if !utils.IsValidField(requestBody.Field) {
		RespondWithError(w, http.StatusBadRequest, "Field must have 2 to 50 characters - and valid characters-words")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	AdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.CleanOldUnacceptedArticlesByDateAndField(dbaccess, dateToCompare, requestBody.Field, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to delete articles: "+err.Error())
		return
	}

	// Build the response
	response := "old unaccepted articles of" + requestBody.Field + " field " + "cleaned successfully until" + requestBody.DateToCompare

	// Respond with the success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, AdminPkey)
}

func CleanAllOldUnacceptedArticlesHandler(w http.ResponseWriter, r *http.Request) {
	// Check the header
	encryptedHeader := r.Header.Get("Encrypted")
	if encryptedHeader != "true" {
		RespondWithError(w, http.StatusBadRequest, "Invalid header. Encrypted must be true!")
		return
	}

	var requestBody struct {
		AdminID uuid.UUID `json:"AdminID"`
	}

	// Decode the request body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if the root admin ID in the URL matches the ID in the request body
	rootAdminIDFromURL := chi.URLParam(r, "rootAdminID")
	if rootAdminIDFromURL != requestBody.AdminID.String() {
		RespondWithError(w, http.StatusBadRequest, "Root admin IDs do not match")
		return
	}

	// Access the GORM.DB connection using dbAccessor
	dbaccess := dbAccessor

	AdminPkey, err := db.GetAdminPublicKey(dbaccess, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = db.CleanAllOldUnacceptedArticles(dbaccess, requestBody.AdminID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to delete articles: "+err.Error())
		return
	}

	response := "all old unaccepted articles cleaned successfully"

	// Respond with the success message
	RespondWithEncryptedJSON(w, http.StatusOK, response, AdminPkey)
}
