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

	router.Patch("/approve-article/{adminID}", adminActionHandler(secretKey, handlers.ApproveArticleHandler))
	router.Patch("/approve-admin/{adminID}", adminActionHandler(secretKey, handlers.ApproveAdminHandler))
	router.Patch("/disapprove-admin/{adminID}", adminActionHandler(secretKey, handlers.DisapproveAdminHandler))
	router.Patch("/mod-permission-admin/{adminID}", adminActionHandler(secretKey, handlers.ModifyAdminPermissionsHandler))

	router.Delete("/delete-article/{adminID}", adminActionHandler(secretKey, handlers.DeleteArticleHandler))
	return router
}

// helper function
func adminActionHandler(secretKey string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminID := chi.URLParam(r, "adminID")
		auth.ActionsMiddleware(adminID, secretKey, handler)(w, r)
	}
}
