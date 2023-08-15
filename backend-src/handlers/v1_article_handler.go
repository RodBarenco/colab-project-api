package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetLatesThousandtArticlesHandler(w http.ResponseWriter, r *http.Request) {
	GetArticlesResponseHandler(w, r, db.GetLatestThousandArticles)
}

func GetLatestFiftyArticlesHandler(w http.ResponseWriter, r *http.Request) {
	GetArticlesResponseHandler(w, r, db.GetLatestFiftyArticles)
}

func GetArticlesByTitleHandler(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	SearchArticlesHandler(w, r, title, db.GetArticlesByTitle)
}

func GetArticlesBySubjectHandler(w http.ResponseWriter, r *http.Request) {
	subject := chi.URLParam(r, "subject")
	SearchArticlesHandler(w, r, subject, db.GetArticlesBySubject)
}

func GetArticlesByFieldHandler(w http.ResponseWriter, r *http.Request) {
	field := chi.URLParam(r, "field")
	SearchArticlesHandler(w, r, field, db.GetArticlesByField)
}

func GetArticlesByAuthorHandler(w http.ResponseWriter, r *http.Request) {
	authorFullName := chi.URLParam(r, "author")

	names := strings.Split(authorFullName, " ")
	if len(names) != 2 {
		RespondWithError(w, http.StatusBadRequest, "Invalid author name format")
		return
	}

	firstName := names[0]
	lastName := names[1]

	dbAccess := dbAccessor

	authorID, err := utils.GetAuthorIDByName(dbAccess, firstName, lastName)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Author not found")
		return
	}

	articles, err := db.GetArticlesByAuthor(dbAccess, authorID) // HERE CALLS the  FUNCTION
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Couldn't get the articles.")
		return
	}

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

		response := res.ArticleResponse{
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
		}

		// Only append if the article is accepted
		if article.IsAccepted {
			articleResponses = append(articleResponses, response)
		}
	}

	// If no accepted articles were found, return a friendly response
	if len(articleResponses) == 0 {
		message := "No article was found."
		RespondWithJSON(w, http.StatusOK, message)
		return
	}

	// Respond with the list of accepted articles
	RespondWithJSON(w, http.StatusOK, articleResponses)
}

func GetArticlesByKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	keywordsQuery := chi.URLParam(r, "keywords") // Obtém o parâmetro da URL (palavras-chave)

	keywords := strings.Split(keywordsQuery, ",") // Divide as palavras-chave separadas por vírgula

	if len(keywords) == 0 {
		RespondWithError(w, http.StatusBadRequest, "No keywords provided")
		return
	}

	dbAccess := dbAccessor // Suponho que você já tenha uma variável "dbAccessor" para acessar o banco de dados

	articles, err := db.GetArticlesByKeywords(dbAccess, keywords...) // Chama a função GetArticlesByKeywords
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error fetching articles")
		return
	}

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

		response := res.ArticleResponse{
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
		}

		// Only append if the article is accepted
		if article.IsAccepted {
			articleResponses = append(articleResponses, response)
		}
	}

	// If no accepted articles were found, return a friendly response
	if len(articleResponses) == 0 {
		message := "No articles was found."
		RespondWithJSON(w, http.StatusOK, message)
		return
	}

	// Respond with the list of accepted articles
	RespondWithJSON(w, http.StatusOK, articleResponses)
}

func GetArticleByIdHandler(w http.ResponseWriter, r *http.Request) {
	articleIDParam := chi.URLParam(r, "id")
	articleID, err := uuid.Parse(articleIDParam)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	dbAccess := dbAccessor

	article, err := db.GetArticleById(dbAccess, articleID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Article not found")
		return
	}

	// Obter os nomes das pessoas que curtiram o artigo
	likedByNames, err := GetLikedByUserNames(dbAccess, articleID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error fetching liked by user names")
		return
	}

	// Obter a quantidade de curtidas
	numLikes := len(likedByNames)

	// Criar a estrutura de resposta
	response := res.ArticleWithLikesResponse{
		Article: article,
		RelatedTables: res.LikesInfo{
			NumLikes:     numLikes,
			LikedByNames: likedByNames,
		},
		Message: "Article and related data retrieved successfully",
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func GetLikedByUsersHandler(w http.ResponseWriter, r *http.Request) {
	articleIDParam := chi.URLParam(r, "id")
	articleID, err := uuid.Parse(articleIDParam)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	dbAccess := dbAccessor

	likedByUsers, err := db.GetLikedByUsers(dbAccess, articleID)
	if err != nil {
		if err.Error() == "Article not found" {
			RespondWithJSON(w, http.StatusNotFound, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Error fetching liked users")
		return
	}

	var likedByUsersRes []res.LikedByUser

	for _, user := range likedByUsers {
		r := res.LikedByUser{
			ID:       user.ID,
			Username: user.FirstName + " " + user.LastName,
		}
		likedByUsersRes = append(likedByUsersRes, r)
	}

	response := res.LikedByUsersResponse{
		LikedByUsers: likedByUsersRes,
		Message:      "Liked users fetched successfully",
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func GetCitingArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articleIDParam := chi.URLParam(r, "id")
	articleID, err := uuid.Parse(articleIDParam)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	dbAccess := dbAccessor

	citingArticles, err := db.GetCitingArticles(dbAccess, articleID)
	if err != nil {
		if errors.Is(err, errors.New("Article not found")) {
			RespondWithJSON(w, http.StatusNotFound, err.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching citing articles")
		}
		return
	}

	var response []res.ArticleCitingCitedRes
	for _, article := range citingArticles {
		r := res.ArticleCitingCitedRes{
			ID:      article.ID,
			Title:   article.Title,
			Message: "Article citing information fetched successfully",
		}
		response = append(response, r)
	}

	RespondWithJSON(w, http.StatusOK, response)
}

func GetCitedByArticlesHandler(w http.ResponseWriter, r *http.Request) {
	articleIDParam := chi.URLParam(r, "id")
	articleID, err := uuid.Parse(articleIDParam)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	dbAccess := dbAccessor

	citedByArticles, err := db.ArticleCitedBy(dbAccess, articleID)
	if err != nil {
		if errors.Is(err, errors.New("Article not found")) {
			RespondWithJSON(w, http.StatusNotFound, err.Error())
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching cited-by articles")
		}
		return
	}

	var response []res.ArticleCitingCitedRes
	for _, article := range citedByArticles {
		r := res.ArticleCitingCitedRes{
			ID:      article.ID,
			Title:   article.Title,
			Message: "Cited-by article information fetched successfully",
		}
		response = append(response, r)
	}

	RespondWithJSON(w, http.StatusOK, response)
}

/// HELPER-------------------------------------------------------------------------------------------------

func GetArticlesResponseHandler(w http.ResponseWriter, r *http.Request, articleGetter func(*gorm.DB) ([]db.Article, error)) {
	dbAccess := dbAccessor

	articles, err := articleGetter(dbAccess)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error fetching articles")
		return
	}

	var articleResponses []res.ArticleResponse

	for _, article := range articles {
		authorName, err := utils.GetAuthorName(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching author name")
			return
		}

		// Get the names of users who liked the article
		likedByNames, err := GetLikedByUserNames(dbAccess, article.ID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching liked by users")
			return
		}

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        likedByNames, // Use the names of users who liked the article
			Shares:         article.Shares,
			CoverImage:     article.CoverImage,
		}

		// Only append if the article is accepted
		if article.IsAccepted {
			articleResponses = append(articleResponses, response)
		}
	}

	// If no accepted articles were found, return a friendly response
	if len(articleResponses) == 0 {
		message := "No accepted articles found."
		RespondWithJSON(w, http.StatusOK, message)
		return
	}

	// Respond with the list of accepted articles
	RespondWithJSON(w, http.StatusOK, articleResponses)
}

//------------

func SearchArticlesHandler(w http.ResponseWriter, r *http.Request, searchParam string, articleGetter func(*gorm.DB, string) ([]db.Article, error)) {
	dbAccess := dbAccessor

	if searchParam == "" || !utils.IsValidArticleSearch(searchParam) {
		RespondWithError(w, http.StatusBadRequest, "Invalid search parameter")
		return
	}

	articles, err := articleGetter(dbAccess, searchParam)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Articles not found")
		return
	}

	var articleResponses []res.ArticleResponse
	for _, article := range articles {
		authorName, err := utils.GetAuthorName(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching author name")
			return
		}

		// Get the names of users who liked the article
		likedByNames, err := GetLikedByUserNames(dbAccess, article.ID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching liked by users")
			return
		}

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        likedByNames, // Use the names of users who liked the article
			Shares:         article.Shares,
			CoverImage:     article.CoverImage,
		}

		// Only append if the article is accepted
		if article.IsAccepted {
			articleResponses = append(articleResponses, response)
		}
	}

	// If no accepted articles were found, return a friendly response
	if len(articleResponses) == 0 {
		message := "No articles found."
		RespondWithJSON(w, http.StatusOK, message)
		return
	}

	RespondWithJSON(w, http.StatusOK, articleResponses)
}

// Function to get name of users that liked a Article !!!!!!!!!!!!!!!!!!
func GetLikedByUserNames(db *gorm.DB, articleID uuid.UUID) ([]string, error) {
	likedByUsers, err := utils.GetNamesOfUsersThatLikedArticles(db, articleID)
	if err != nil {
		return nil, err
	}

	likedByNames := make([]string, len(likedByUsers))
	for i, user := range likedByUsers {
		likedByNames[i] = user.FirstName + " " + user.LastName
	}

	return likedByNames, nil
}
