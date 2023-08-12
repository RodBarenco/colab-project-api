package routes

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func URLValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !utils.ValidateURL(path) {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		// Se a URL for válida, continue para o próximo handler
		next.ServeHTTP(w, r)
	})
}

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

	router.Use(URLValidationMiddleware)

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
