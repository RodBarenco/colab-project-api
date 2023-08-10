package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/RodBarenco/colab-project-api/utils"
)

func CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	var newArticle db.ArticleParams
	err := json.NewDecoder(r.Body).Decode(&newArticle)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check for empty fields in 'newArticle'
	if newArticle.Title == "" || !utils.IsValidArticleTitle(newArticle.Title) {
		RespondWithError(w, http.StatusBadRequest, "Invalid or missing article title")
		return
	}

	if newArticle.Subject == "" || !utils.IsValidArticleSubject(newArticle.Subject) {
		RespondWithError(w, http.StatusBadRequest, "Invalid or missing article subject")
		return
	}

	if newArticle.Field == "" || !utils.IsValidArticleField(newArticle.Field) {
		RespondWithError(w, http.StatusBadRequest, "Invalid or missing article field")
		return
	}

	if len(newArticle.File) == 0 || !utils.IsValidArticleFile(newArticle.File) {
		RespondWithError(w, http.StatusBadRequest, "Invalid or missing article file")
		return
	}

	if !utils.IsValidArticleDescription(newArticle.Description) {
		RespondWithError(w, http.StatusBadRequest, "Invalid article description")
		return
	}

	if !utils.IsValidArticleKeywords(newArticle.Keywords) {
		RespondWithError(w, http.StatusBadRequest, "Invalid article keywords")
		return
	}

	if newArticle.CoAuthors != "" && !utils.IsValidArticleCoAuthors(newArticle.CoAuthors) {
		RespondWithError(w, http.StatusBadRequest, "Invalid article co-authors")
		return
	}

	if newArticle.CoverImage != "" && !utils.IsValidArticleCoverImage(newArticle.CoverImage) {
		RespondWithError(w, http.StatusBadRequest, "Invalid article cover image")
		return
	}

	// Access the gorm.DB connection using dbAccessor
	dbAccess := dbAccessor

	// Call the CreateArticle function passing the ArticleParams object and the database connection
	err = db.CreateArticle(dbAccess, newArticle)
	if err != nil {
		errorMessage := fmt.Sprintf("Error creating article: %v", err)
		RespondWithError(w, http.StatusInternalServerError, errorMessage)
		return
	}

	response := res.CreateArticleRes{
		Article: res.ArticleCreatedResponse{
			Title:   newArticle.Title,
			Subject: newArticle.Subject,
			Field:   newArticle.Field,
		},
		Message: "Article created successfully!",
	}

	RespondWithJSON(w, http.StatusCreated, response)
}
