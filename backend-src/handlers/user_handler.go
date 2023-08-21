package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func CreateArticleHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
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

	RespondToLoggedInUser(w, r, encryptResponse, response, newArticle.AuthorID)
}

func GetRecommendedArticlesHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	userIDString := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid userID")
		return
	}

	monthsAgoString := r.URL.Query().Get("monthsAgo")
	var monthsAgo int
	if monthsAgoString != "" {
		monthsAgo, err = strconv.Atoi(monthsAgoString)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid monthsAgo")
			return
		}
	}

	dbAccess := dbAccessor

	articles, otherArticles, err := db.GetRecommendedArticles(dbAccess, userID, monthsAgo)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to get recommended articles")
		return
	}

	// Create ArticleResponse slices from articles and otherArticles
	var articleResponses []res.ArticleResponse
	for _, article := range articles {
		authorName, err := utils.GetAuthorName(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching author name")
			return
		}

		likedByNames, err := GetLikedByUserNames(dbAccess, article.ID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching liked by users")
			return
		}

		articleResponses = append(articleResponses, res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        likedByNames,
			Shares:         article.Shares,
			CoverImage:     article.CoverImage,
		})
	}

	// Create ArticleResponse slices for otherArticles
	var otherArticleResponses []res.ArticleResponse
	for _, article := range otherArticles {
		authorName, err := utils.GetAuthorName(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching author name")
			return
		}

		likedByNames, err := GetLikedByUserNames(dbAccess, article.ID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching liked by users")
			return
		}

		otherArticleResponses = append(otherArticleResponses, res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        likedByNames,
			Shares:         article.Shares,
			CoverImage:     article.CoverImage,
		})
	}

	response := struct {
		Articles      []res.ArticleResponse `json:"articles"`
		OtherArticles []res.ArticleResponse `json:"other_articles"`
	}{
		Articles:      articleResponses,
		OtherArticles: otherArticleResponses,
	}

	RespondToLoggedInUser(w, r, encryptResponse, response, userID)
}

func AddPublicKeyHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	var requestData db.AddPublicKeyRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userIDFromURL := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDFromURL)
	if userID != requestData.UserID {
		RespondWithError(w, http.StatusBadRequest, "Mismatched user IDs")
		return
	}

	// Verificar formato Base64 da chave pública
	if _, err := base64.StdEncoding.DecodeString(requestData.PublicKey); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid public key format")
		return
	}

	err = db.AddPublicKeyToUser(dbAccessor, userID, requestData.PublicKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error adding public key")
		return
	}

	var message = "message :  Public key added successfully"

	RespondToLoggedInUser(w, r, encryptResponse, message, userID)
}

//--------------	HELPER ----------------//

func RespondToLoggedInUser(w http.ResponseWriter, r *http.Request, encryptResponse bool, response interface{}, userID uuid.UUID) {
	// Verificar se a resposta deve ser encriptada
	if encryptResponse {
		encryptedHeader := r.Header.Get("Encrypted")
		if encryptedHeader != "true" {
			RespondWithError(w, http.StatusInternalServerError, "Request must be encrypted")
			return
		}

		// Obter a chave pública do usuário
		pkey, err := db.GetUserPublicKey(dbAccessor, userID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error getting user public key")
			return
		}

		// Encriptar e enviar a resposta
		RespondWithEncryptedJSON(w, http.StatusOK, response, pkey)
	} else {
		// Enviar a resposta sem encriptação
		RespondWithJSON(w, http.StatusOK, response)
	}
}
