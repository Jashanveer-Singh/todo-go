package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Jashanveer-Singh/todo-go/internal/ports"
)

func NewAuthMiddleware(tokenProvider ports.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{
		tokenProvider: tokenProvider,
	}
}

type AuthMiddleware struct {
	tokenProvider ports.TokenProvider
}

func getBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("Invalid authorization header format")
	}
	return authHeader[7:], nil
}

func (am AuthMiddleware) isAuthenticatedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getBearerToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		claims, err := am.tokenProvider.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "claims", claims))
		next.ServeHTTP(w, r)
		// return claims, nil
	}
}
