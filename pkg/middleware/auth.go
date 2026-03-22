package middleware

import (
	"context"
	"net/http"
	"strings"

	"example.com/ecommerce/internal/auth"
)

type contextKey string

const UserIdKey contextKey = "user_id"
const EmailKey contextKey = "email"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// inject user info into context so handlers can read it
		ctx := context.WithValue(r.Context(), UserIdKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
