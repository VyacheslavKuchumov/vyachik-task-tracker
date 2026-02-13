package tracker

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler, userStore types.UserStore) {
	r.Route("/goals", func(r chi.Router) {
		r.Get("/", auth.WithJWTAuth(handler.HandleGetGoals, userStore))
		r.Post("/", auth.WithJWTAuth(handler.HandleCreateGoal, userStore))
		r.Post("/{goalID}/tasks", auth.WithJWTAuth(handler.HandleCreateTask, userStore))
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/assigned", auth.WithJWTAuth(handler.HandleGetAssignedTasks, userStore))
		r.Put("/{taskID}/assign", auth.WithJWTAuth(handler.HandleAssignTask, userStore))
	})
}
