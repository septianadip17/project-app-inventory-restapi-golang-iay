package middleware

import (
	"context"
	"net/http"
	"project-app-inventory-restapi-golang-iay/internal/entity"
	"strings"
	"time"
	// import entity & response pkg
)

// AuthMiddleware memvalidasi token dari DB
func (m *MiddlewareManager) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Return 401 Unauthorized
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Cek ke DB via Repository
		session, err := m.sessionRepo.GetSessionByToken(r.Context(), tokenStr)
		if err != nil || session.RevokedAt != nil || session.ExpiredAt.Before(time.Now()) {
			// Return 401 Unauthorized (Token invalid/expired/revoked)
			return
		}

		// Ambil data user untuk context role
		user, _ := m.userRepo.GetByID(r.Context(), session.UserID)

		// Simpan user info ke Context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleMiddleware membatasi akses berdasarkan role
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value("user").(*entity.User)
			if !ok {
				// 401 Unauthorized
				return
			}

			for _, role := range allowedRoles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			// 403 Forbidden
			http.Error(w, "Forbidden Access", http.StatusForbidden)
		})
	}
}
