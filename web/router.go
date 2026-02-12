package web

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler, userStore types.UserStore) {
	r.Handle("/static/*", handler.StaticHandler())

	r.Get("/", handler.HandleHome)
	r.Get("/login", handler.HandleLoginPage)
	r.Get("/register", handler.HandleRegisterPage)
	r.Get("/goals", auth.WithJWTPageAuth(handler.HandleGoalsPage, userStore))
	r.Get("/goals/edit", auth.WithJWTPageAuth(handler.HandleGoalEditPage, userStore))
	r.Get("/tasks", auth.WithJWTPageAuth(handler.HandleTasksPage, userStore))
	r.Get("/tasks/edit", auth.WithJWTPageAuth(handler.HandleTaskEditPage, userStore))

	r.Route("/htmx", func(r chi.Router) {
		r.Get("/goals", auth.WithJWTAuth(handler.HandleHTMXGoals, userStore))
		r.Get("/goals/card", auth.WithJWTAuth(handler.HandleHTMXGoalCard, userStore))
		r.Get("/goals/card/{goalID}", auth.WithJWTAuth(handler.HandleHTMXGoalCard, userStore))
		r.Post("/goals/save", auth.WithJWTAuth(handler.HandleHTMXGoalSave, userStore))
		r.Get("/tasks", auth.WithJWTAuth(handler.HandleHTMXTasks, userStore))
		r.Get("/tasks/card", auth.WithJWTAuth(handler.HandleHTMXTaskCard, userStore))
		r.Get("/tasks/card/{taskID}", auth.WithJWTAuth(handler.HandleHTMXTaskCard, userStore))
		r.Post("/tasks/save", auth.WithJWTAuth(handler.HandleHTMXTaskSave, userStore))
	})
}
