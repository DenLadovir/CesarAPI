package middleware

import (
	"CesarAPI/utils"
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const usernameKey contextKey = "username"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получение токена из заголовка Authorization
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Проверка формата токена (Bearer <токен>)
		if !strings.HasPrefix(tokenString, "Bearer ") {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Валидация токена
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			log.Printf("JWT validation error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Добавление данных в контекст
		ctx := context.WithValue(r.Context(), usernameKey, claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
