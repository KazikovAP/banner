package middlewares

import (
	response "banner/internal/lib/api/responses"
	jwt "banner/internal/lib/auth/jwt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

func TokenAuthMiddleware(jwtManager *jwt.JWTSecret, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		token := strings.Fields(tokenString)
		if len(token) != 2 || token[0] != "Bearer" {
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		_, err := jwtManager.VerifyToken(token[1])
		if err != nil {
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TokenAuthAndRoleMiddleware(jwtManager *jwt.JWTSecret, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		token := strings.Fields(tokenString)
		if len(token) != 2 || token[0] != "Bearer" {
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		claims, err := jwtManager.VerifyToken(token[1])
		if err != nil {
			render.JSON(w, r, response.Error("Unauthorized"))
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			render.JSON(w, r, response.Error("Forbidden"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
