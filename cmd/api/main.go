func main() {
	// ... Load Config, DB, Logger ...

	r := chi.NewRouter()
	r.Use(middleware.Logger) // Logger basic

	// Public Routes
	r.Post("/login", authHandler.Login)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(mw.AuthMiddleware) // Pasang middleware Auth

		// Rute untuk STAFF, ADMIN, SUPER_ADMIN
		r.Get("/items", itemHandler.GetAllItems)
		r.Get("/items/minimum-stock", itemHandler.GetLowStockItems) // Fitur Cek Stok Minimum

		// Rute hanya untuk ADMIN & SUPER_ADMIN
		r.Group(func(r chi.Router) {
			r.Use(mw.RoleMiddleware("admin", "super_admin"))

			r.Post("/items", itemHandler.CreateItem)
			r.Put("/items/{id}", itemHandler.UpdateItem)
			r.Delete("/items/{id}", itemHandler.DeleteItem)

			r.Post("/categories", categoryHandler.CreateCategory)
		})

		// Rute Khusus SUPER_ADMIN
		r.Group(func(r chi.Router) {
			r.Use(mw.RoleMiddleware("super_admin"))

			r.Post("/users", userHandler.RegisterUser) // Create user baru
			r.Delete("/users/{id}", userHandler.DeleteUser)
		})
	})

	// Start Server
}