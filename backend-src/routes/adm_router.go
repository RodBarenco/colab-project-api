package routes

import (
	"net/http"

	"github.com/RodBarenco/colab-project-api/auth"
	"github.com/RodBarenco/colab-project-api/handlers"
	"github.com/go-chi/chi"
)

func AdmRoutes(secretKey string) http.Handler {
	router := chi.NewRouter()

	router.Get("/testreadiness", handlers.HandlerReadiness)
	router.Get("/testerror", handlers.HandlerError)

	router.Get("/unaccepted-article-id/{adminID}", adminActionHandler(secretKey, handlers.GetUnacceptedArticlesIDHandler))
	router.Get("/unaccepted-article-by-field/{adminID}", adminActionHandler(secretKey, handlers.GetUnacceptedArticlesByFieldHandler))

	router.Patch("/approve-article/{adminID}", adminActionHandler(secretKey, handlers.ApproveArticleHandler))
	router.Patch("/approve-admin/{adminID}", adminActionHandler(secretKey, handlers.ApproveAdminHandler))
	router.Patch("/disapprove-admin/{adminID}", adminActionHandler(secretKey, handlers.DisapproveAdminHandler))
	router.Patch("/mod-permission-admin/{adminID}", adminActionHandler(secretKey, handlers.ModifyAdminPermissionsHandler))

	router.Delete("/delete-article/{adminID}", adminActionHandler(secretKey, handlers.DeleteArticleHandler))
	router.Delete("/delete-user/{adminID}", adminActionHandler(secretKey, handlers.DeleteUserHandler))
	router.Delete("/delete-admin/{adminID}", adminActionHandler(secretKey, handlers.DeleteAdminHandler))
	router.Delete("/clean-articles/{adminID}", adminActionHandler(secretKey, handlers.CleanAllOldUnacceptedArticlesHandler))
	router.Delete("/clean-articles-by-date/{adminID}", adminActionHandler(secretKey, handlers.CleanOldUnacceptedArticlesByDateHandler))
	router.Delete("/clean-articles-by-date-and-field/{adminID}", adminActionHandler(secretKey, handlers.CleanOldUnacceptedArticlesByDateAndFieldHandler))

	return router
}

// helper function
func adminActionHandler(secretKey string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminID := chi.URLParam(r, "adminID")
		auth.ActionsMiddleware(adminID, secretKey, handler)(w, r)
	}
}
