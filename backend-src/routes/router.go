package routes

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/rsakeys"
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

	pkey, err := rsakeys.ReadPrivateKeyFromFile()
	if err != nil {
		panic("no private key avalable!")
	}

	// Roteador para usuários não logados
	v1Router := GeneralRoutes()
	router.Mount("/v1", auth.DecryptionMiddleware(v1Router, pkey))

	// Roteador para usuários logados
	v2Router := UserRoutes(secretKey)
	router.Mount("/v2", auth.DecryptionMiddleware(auth.AuthMiddleware("user", secretKey, v2Router), pkey))

	// Roteador para adms logados
	v3Router := AdmRoutes(secretKey)
	router.Mount("/v3", auth.DecryptionMiddleware(auth.AuthMiddleware("admin", secretKey, v3Router), pkey))

	return router
}
