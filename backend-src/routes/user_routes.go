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
	router.Get("/testreadiness/{userID}", userActionHandler(secretKey, handlers.HandlerReadinessEncrypted, false))

	router.Get("/home-articles/{userID}", userActionHandler(secretKey, handlers.GetRecommendedArticlesHandler, false))
	router.Get("/is-article-liked/{userID}", userActionHandler(secretKey, handlers.IsArticleLikedByUserHandler, false))

	router.Post("/create-article/{userID}", userActionHandler(secretKey, handlers.CreateArticleHandler, false))
	router.Patch("/like-article/{userID}", userActionHandler(secretKey, handlers.AddUserToLikedByFromArticleHandler, false))
	router.Patch("/unlike-article/{userID}", userActionHandler(secretKey, handlers.RemoveUserFromLikedByFromArticleHandler, false))
	router.Patch("/add-cited-article/{userID}", userActionHandler(secretKey, handlers.AddCitationHandler, false))
	router.Patch("/remove-cited-article/{userID}", userActionHandler(secretKey, handlers.RemoveCitationHandler, false))
	router.Patch("/add-key/{userID}", userActionHandler(secretKey, handlers.AddPublicKeyHandler, false))
	router.Patch("/follow/{userID}", userActionHandler(secretKey, handlers.FollowUserHandler, false))
	router.Patch("/unfollow/{userID}", userActionHandler(secretKey, handlers.UnfollowUserHandler, false))

	//encypted responses----------------------------------------------------------------------------

	router.Get("/testreadiness/{userID}", userActionHandler(secretKey, handlers.HandlerReadinessEncrypted, true))

	router.Get("/home-articles-ecpt/{userID}", userActionHandler(secretKey, handlers.GetRecommendedArticlesHandler, true))
	router.Get("/is-article-liked-ecpt/{userID}", userActionHandler(secretKey, handlers.IsArticleLikedByUserHandler, true))

	router.Post("/create-article-ecpt/{userID}", userActionHandler(secretKey, handlers.CreateArticleHandler, true))
	router.Patch("/like-article-ecpt/{userID}", userActionHandler(secretKey, handlers.AddUserToLikedByFromArticleHandler, true))
	router.Patch("/unlike-article-ecpt/{userID}", userActionHandler(secretKey, handlers.RemoveUserFromLikedByFromArticleHandler, true))
	router.Patch("/add-cited-article-ecpt/{userID}", userActionHandler(secretKey, handlers.AddCitationHandler, true))
	router.Patch("/remove-cited-article-ecpt/{userID}", userActionHandler(secretKey, handlers.RemoveCitationHandler, true))
	router.Patch("/add-key-ecpt/{userID}", userActionHandler(secretKey, handlers.AddPublicKeyHandler, true))
	router.Patch("/follow-ecpt/{userID}", userActionHandler(secretKey, handlers.FollowUserHandler, true))
	router.Patch("/unfollow-ecpt/{userID}", userActionHandler(secretKey, handlers.UnfollowUserHandler, true))

	return router
}

// HELPER FUNCTION
func userActionHandler(secretKey string, handler func(http.ResponseWriter, *http.Request, bool), e bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		auth.ActionsMiddleware(userID, secretKey, handlerFuncWithEncryption(handler, e)).ServeHTTP(w, r)
	}
}

func handlerFuncWithEncryption(handler func(http.ResponseWriter, *http.Request, bool), e bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, e)
	}
}
