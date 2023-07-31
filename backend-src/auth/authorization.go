package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func Authorize(tokenString string, secret string, requiredRole string) (int, error) {
	// Verify and parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method used in the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("failed to verify the JWT token: %v", err)
	}

	// Check if the token is valid and not expired
	if !token.Valid {
		return http.StatusUnauthorized, fmt.Errorf("invalid or expired JWT token")
	}

	// Verify the "aud" claim in the token to ensure the user has the correct role
	if requiredRole != "" {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return http.StatusUnauthorized, fmt.Errorf("failed to get the token claims")
		}

		aud, ok := claims["aud"].(string)
		if !ok {
			return http.StatusUnauthorized, fmt.Errorf("failed to get the 'aud' claim from the token")
		}

		if aud != requiredRole {
			return http.StatusForbidden, fmt.Errorf("user is not authorized for this operation")
		}
	}

	return http.StatusOK, nil
}
