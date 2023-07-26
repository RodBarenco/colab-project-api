package routes

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/handlers"
	"github.com/go-chi/chi"
)

func UserRoutes() http.Handler {
	router := chi.NewRouter()

	router.Get("/testreadiness", handlers.HandlerReadiness)
	router.Get("/testerror", handlers.HandlerError)

	return router
}
