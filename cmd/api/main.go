package main

import (
	"net/http"
	"project-app-inventory-restapi-golang-iay/internal/config"
	"project-app-inventory-restapi-golang-iay/internal/handler"
	"project-app-inventory-restapi-golang-iay/internal/middleware"
	"project-app-inventory-restapi-golang-iay/internal/repository"
	"project-app-inventory-restapi-golang-iay/internal/service"
	"project-app-inventory-restapi-golang-iay/pkg/database"
	"project-app-inventory-restapi-golang-iay/pkg/logger"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, _ := config.LoadConfig()
	log := logger.NewLogger()
	db, _ := database.NewPostgresDB(cfg.DBUrl)
	defer db.Close()

	// Init Layers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Routes
	r.Post("/login", authHandler.Login)

	// Protected Routes (Contoh Skeleton)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.VerifyToken)

		// Contoh endpoint dummy untuk test role
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You are logged in!"))
		})

		// Routing untuk Super Admin only
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleCheck("super_admin"))
			r.Get("/admin-only", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello Super Admin"))
			})
		})
	})

	log.Info("Server running on " + cfg.ServerPort)
	http.ListenAndServe(cfg.ServerPort, r)
}
