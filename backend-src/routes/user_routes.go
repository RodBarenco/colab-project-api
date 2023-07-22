package routes

import (
	"net/http"

	"github.com/go-chi/chi"
)

func UserRoutes() http.Handler {
	router := chi.NewRouter()

	// Definir rotas para os usuários logados aqui

	return router
}
