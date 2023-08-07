package auth

import (
	"net/http"
	"strings"
)

func AuthMiddleware(requiredRole string, secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the request header
		tokenString := extractTokenFromHeader(r)

		// If the token is not present in the header, try getting it from the request context
		if tokenString == "" {
			token, ok := r.Context().Value("jwtToken").(string)
			if !ok {
				http.Error(w, "JWT token not found in the request", http.StatusUnauthorized)
				return
			}
			tokenString = token
		}

		// Call the Authorize function for token validation
		statusCode, err := Authorize(tokenString, secret, requiredRole)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// Call the next handler with the updated context.
		next.ServeHTTP(w, r)
	})
}

func ActionsMiddleware(id, secret string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the request header
		tokenString := extractTokenFromHeader(r)

		// If the token is not present in the header, try getting it from the request context
		if tokenString == "" {
			token, ok := r.Context().Value("jwtToken").(string)
			if !ok {
				http.Error(w, "JWT token not found in the request", http.StatusUnauthorized)
				return
			}
			tokenString = token
		}

		// Call the Authorize function for token validation
		statusCode, err := AuthorizeActions(tokenString, secret, id)
		if err != nil {
			http.Error(w, err.Error(), statusCode)
			return
		}

		// If the authorization is successful, call the handler function
		handler(w, r)
	}
}

// Helper function to extract the JWT token from the Authorization header
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Expecting the Authorization header value to be in the format "Bearer <token>"
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) == 2 && splitToken[0] == "Bearer" {
			return splitToken[1]
		}
	}
	return ""
}
