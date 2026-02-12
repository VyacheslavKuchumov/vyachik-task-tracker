package user

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Post("/auth/login", handler.HandleWebLogin)
	r.Post("/auth/register", handler.HandleWebRegister)
	r.Post("/auth/logout", handler.HandleLogout)

	r.Post("/api/v1/login", handler.HandleLogin)
	r.Post("/api/v1/register", handler.HandleRegister)
}
