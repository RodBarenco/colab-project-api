package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func MainRouter() http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Roteador para usuários não logados
	v1Router := GeneralRoutes()
	router.Mount("/v1", v1Router)

	// Roteador para usuários logados
	v2Router := UserRoutes()
	router.Mount("/v2", v2Router)

	return router
}
