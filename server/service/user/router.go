package user

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Post("/login", handler.HandleLogin)
	r.Post("/register", handler.HandleRegister)
	r.Get("/profile", handler.HandleGetProfile)
	r.Put("/profile", handler.HandleUpdateProfile)
	r.Put("/profile/password", handler.HandleUpdatePassword)
	r.Get("/users/lookup", handler.HandleListUsers)
}
