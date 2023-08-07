package routes

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func MainRouter(secretKey string) http.Handler {
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
	v2Router := UserRoutes(secretKey)
	router.Mount("/v2", auth.AuthMiddleware("user", secretKey, v2Router))

	// Roteador para adms logados
	v3Router := AdmRoutes(secretKey)
	router.Mount("/v3", auth.AuthMiddleware("admin", secretKey, v3Router))

	return router
}
