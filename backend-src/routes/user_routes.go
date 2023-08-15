package routes

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/handlers"
	"github.com/go-chi/chi"
)

func UserRoutes(secretKey string) http.Handler {
	router := chi.NewRouter()

	router.Get("/testerror", handlers.HandlerReadiness)
	router.Get("/testerror", handlers.HandlerError)
	router.Get("/testreadiness/{userID}", userActionHandler(secretKey, handlers.HandlerReadiness))

	router.Get("/home-articles/{userID}", userActionHandler(secretKey, handlers.GetRecommendedArticlesHandler))
	router.Get("/is-article-liked/{userID}", userActionHandler(secretKey, handlers.IsArticleLikedByUserHandler))

	router.Post("/create-article/{userID}", userActionHandler(secretKey, handlers.CreateArticleHandler))
	router.Patch("/like-article/{userID}", userActionHandler(secretKey, handlers.AddUserToLikedByHandler))
	router.Patch("/unlike-article/{userID}", userActionHandler(secretKey, handlers.RemoveUserFromLikedByHandler))
	router.Patch("/add-cited-article/{userID}", userActionHandler(secretKey, handlers.AddCitationHandler))
	router.Patch("/remove-cited-article/{userID}", userActionHandler(secretKey, handlers.RemoveCitationHandler))

	return router
}

// HELPER FUNCTION
func userActionHandler(secretKey string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		auth.ActionsMiddleware(userID, secretKey, handler).ServeHTTP(w, r)
	}
}
