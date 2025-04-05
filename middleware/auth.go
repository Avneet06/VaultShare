package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey")

type contextKey string

const UserEmailKey contextKey = "userEmail"

func JWTAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("JWT parse error:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		email := claims["email"].(string)
		ctx := context.WithValue(r.Context(), UserEmailKey, email)
		next(w, r.WithContext(ctx))
	}
}
