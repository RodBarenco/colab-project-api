package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// -------------------------LIKES ------------------------------------//
func AddUserToLikedByFromArticleHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	var params db.LikeArticleRequestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	userIDFromURL := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDFromURL)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	if userID != params.UserID {
		RespondWithError(w, http.StatusBadRequest, "User ID mismatch")
		return
	}

	if _, err := uuid.Parse(params.ArticleID.String()); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Article ID format")
		return
	}

	dbAccess := dbAccessor

	// Chame a função para adicionar o usuário à lista de "likedBy"
	err = db.AddUserToLikedByFromArticle(dbAccess, params.ArticleID, params.UserID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error adding user to likedBy")
		return
	}

	response := res.CreateLikeArticleResponse{
		LikeArticleResponse: res.LikeArticleResponse{
			UserID:    params.UserID,
			ArticleID: params.ArticleID,
		},
		Message: "Article was liked",
	}

	RespondToLoggedInUser(w, r, encryptResponse, response, userID)
}

func RemoveUserFromLikedByFromArticleHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	var params db.LikeArticleRequestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	userIDFromURL := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDFromURL)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	if userID != params.UserID {
		RespondWithError(w, http.StatusBadRequest, "User ID mismatch")
		return
	}

	if _, err := uuid.Parse(params.ArticleID.String()); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Article ID format")
		return
	}

	dbAccess := dbAccessor

	// Chame a função para remover o usuário da lista de "likedBy"
	err = db.RemoveUserFromLikedByFromArticle(dbAccess, params.ArticleID, params.UserID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error removing user from likedBy")
		return
	}

	response := res.CreateLikeArticleResponse{
		LikeArticleResponse: res.LikeArticleResponse{
			UserID:    params.UserID,
			ArticleID: params.ArticleID,
		},
		Message: "Aricle was unliked",
	}

	RespondToLoggedInUser(w, r, encryptResponse, response, userID)
}

// just if needed by frontend

func IsArticleLikedByUserHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	var params db.LikeArticleRequestParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	userIDFromURL := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDFromURL)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	if userID != params.UserID {
		RespondWithError(w, http.StatusBadRequest, "User ID mismatch")
		return
	}

	if _, err := uuid.Parse(params.ArticleID.String()); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Article ID format")
		return
	}

	dbAccess := dbAccessor

	isLiked, err := db.IsArticleLikedByUser(dbAccess, params.ArticleID, params.UserID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error checking if liked by user")
		return
	}

	response := struct {
		IsLiked bool `json:"is_liked"`
	}{
		IsLiked: isLiked,
	}

	RespondToLoggedInUser(w, r, encryptResponse, response, userID)
}

//-------------------------CITING-----------------------------//

func AddCitationHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	var params db.CitingArticleRequestParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	userIDFromURL := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDFromURL)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	if userID != params.UserID {
		RespondWithError(w, http.StatusBadRequest, "User ID mismatch")
		return
	}

	if params.CitingArticleID == params.CitedArticleID {
		RespondWithError(w, http.StatusBadRequest, "Citing and cited articles must have different IDs")
		return
	}

	dbAccess := dbAccessor

	// Chame a função para adicionar a citação
	err = db.AddCitation(dbAccess, params.CitingArticleID, params.CitedArticleID, params.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "User is not authorized to add citation to this article - or article doesn't exist") {
			RespondWithJSON(w, http.StatusUnauthorized, err.Error())
		} else if strings.Contains(err.Error(), "Cited article not found") {
			RespondWithJSON(w, http.StatusBadRequest, err.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Error adding citation")
		}
		return
	}

	response := res.ArticleCitationResponse{
		CitedArticleID:  params.CitedArticleID.String(),
		CitingArticleID: params.CitingArticleID.String(),
		Message:         "Citation added successfully.",
	}

	RespondToLoggedInUser(w, r, encryptResponse, response, userID)
}

func RemoveCitationHandler(w http.ResponseWriter, r *http.Request, encryptResponse bool) {
	var params db.CitingArticleRequestParams

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	userIDFromURL := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDFromURL)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	if userID != params.UserID {
		RespondWithError(w, http.StatusBadRequest, "User ID mismatch")
		return
	}

	if params.CitingArticleID == params.CitedArticleID {
		RespondWithError(w, http.StatusBadRequest, "Citing and cited articles must have different IDs")
		return
	}

	dbAccess := dbAccessor

	// Chame a função para remover a citação
	err = db.RemoveCitation(dbAccess, params.CitingArticleID, params.CitedArticleID, params.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "User is not authorized to remove citation to this article - or article doesn't exist") {
			RespondWithJSON(w, http.StatusUnauthorized, err.Error())
		} else if strings.Contains(err.Error(), "Cited article not found") {
			RespondWithJSON(w, http.StatusBadRequest, err.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Error removing citation")
		}
		return
	}

	response := res.ArticleCitationResponse{
		CitedArticleID:  params.CitedArticleID.String(),
		CitingArticleID: params.CitingArticleID.String(),
		Message:         "Citation removed successfully.",
	}

	RespondToLoggedInUser(w, r, encryptResponse, response, userID)
}
