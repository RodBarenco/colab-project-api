package routes

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/RodBarenco/colab-project-api/handlers"
)

func GeneralRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/testreadiness", handlers.HandlerReadiness)
	router.Get("/testerror", handlers.HandlerError)

	// articles
	router.Get("/articles", handlers.GetLatesThousandtArticlesHandler)
	router.Get("/articles/home", handlers.GetLatestFiftyArticlesHandler)
	router.Get("/articles/Name/{name}", handlers.GetArticlesByNameHandler)
	router.Get("/articles/title/{title}", handlers.GetArticlesByTitletHandler)
	router.Get("/articles/subject/{subject}", handlers.GetArticlesBySubjectHandler)
	router.Get("/articles/author/{author}", handlers.GetArticlesByAuthorHandler)
	router.Get("/articles/field/{field}", handlers.GetArticlesByFieldHandler)
	router.Get("/articles/keywords/{field}", handlers.GetArticlesByKeywordsHandler)
	router.Get("/articles/id", handlers.GetLatesArticleByIdHandler)

	// art work
	router.Get("/artworks", handlers.GetLatestArtworksHandler)
	router.Get("/artworks/Name/{name}", handlers.GetArtworksByNameHandler)
	router.Get("/artworks/subject/{subject}", handlers.GetArtworksBySubjectHandler)
	router.Get("/artworks/author/{author}", handlers.GetArtworksByAuthorHandler)
	router.Get("/artworks/field/{field}", handlers.GetArtworksByFieldHandler)
	router.Get("/artworks/keywords/{field}", handlers.GetArtworksByKeywordsHandler)
	router.Get("/artworks/id", handlers.GetLatesArtworkByIdHandler)

	// register-login
	router.Post("/register", handlers.RegisterHandler)
	router.Post("/login", handlers.LoginHandler)

	return router
}
