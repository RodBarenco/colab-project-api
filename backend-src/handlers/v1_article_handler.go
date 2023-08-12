package handlers

import (
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

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        article.LikedBy,
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

	// Aqui você pode criar a estrutura de resposta ou retornar os artigos diretamente, dependendo da sua necessidade.
	// Vou criar uma resposta simplificada aqui.
	var articleResponses []res.ArticleResponse

	for _, article := range articles {
		authorName, err := utils.GetAuthorName(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Error fetching author name")
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
			LikedBy:        article.LikedBy,
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

	RespondWithJSON(w, http.StatusOK, article)
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

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        article.LikedBy,
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

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     authorName,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        article.LikedBy,
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
