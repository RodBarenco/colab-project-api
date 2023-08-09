package handlers

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/db"
	"github.com/RodBarenco/colab-project-api/res"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/go-chi/chi"
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
	// TODO
}

func GetArticlesByKeywordsHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func GetLatesArticleByIdHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implemente a lógica para obter o artigo mais recente por ID.
}

/// HELPER-------------------------------------------------------------------------------------------------

func GetArticlesResponseHandler(w http.ResponseWriter, r *http.Request, articleGetter func(*gorm.DB) ([]db.Article, error)) {
	dbAccess := dbAccessor

	articles, err := articleGetter(dbAccess)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao obter os artigos")
		return
	}

	var articleResponses []res.ArticleResponse
	for _, article := range articles {
		user, err := db.GetUserByID(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Erro ao obter nome do autor")
			return
		}
		var author = (user.FirstName + " " + user.LastName)

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     author,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        article.LikedBy,
			Shares:         article.Shares,
			CoverImage:     article.CoverImage,
		}
		articleResponses = append(articleResponses, response)
	}

	RespondWithJSON(w, http.StatusOK, articleResponses)
}

func SearchArticlesHandler(w http.ResponseWriter, r *http.Request, searchParam string, articleGetter func(*gorm.DB, string) ([]db.Article, error)) {
	dbAccess := dbAccessor

	if searchParam == "" || !utils.ArticleSearchIsValid(searchParam) {
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
		user, err := db.GetUserByID(dbAccess, article.AuthorID)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Author not found")
			return
		}
		var author = (user.FirstName + " " + user.LastName)

		response := res.ArticleResponse{
			ID:             article.ID,
			Title:          article.Title,
			AuthorName:     author,
			Subject:        article.Subject,
			Field:          article.Field,
			Description:    article.Description,
			Keywords:       article.Keywords,
			SubmissionDate: article.SubmissionDate,
			LikedBy:        article.LikedBy,
			Shares:         article.Shares,
			CoverImage:     article.CoverImage,
		}
		articleResponses = append(articleResponses, response)
	}

	RespondWithJSON(w, http.StatusOK, articleResponses)
}
