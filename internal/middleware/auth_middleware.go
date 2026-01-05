package middleware

import (
	"context"
	"net/http"
	"project-app-inventory-restapi-golang-iay/internal/entity"
	"project-app-inventory-restapi-golang-iay/internal/repository"
	"project-app-inventory-restapi-golang-iay/pkg/response"
	"strings"
	"time"
)

type AuthMiddleware struct {
	repo *repository.UserRepository
}

func NewAuthMiddleware(repo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{repo: repo}
}

func (m *AuthMiddleware) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		session, err := m.repo.GetSession(r.Context(), tokenStr)
		if err != nil || session.RevokedAt != nil || session.ExpiredAt.Before(time.Now()) {
			response.Error(w, http.StatusUnauthorized, "token expired or invalid")
			return
		}

		user, err := m.repo.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleCheck(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user := r.Context().Value("user").(*entity.User)

			for _, role := range allowedRoles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.Error(w, http.StatusForbidden, "forbidden access")
		})
	}
}
