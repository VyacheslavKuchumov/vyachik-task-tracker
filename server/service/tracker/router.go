package tracker

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/goals", func(r chi.Router) {
		r.Get("/", handler.HandleGetGoals)
		r.Post("/", handler.HandleCreateGoal)
		r.Put("/{goalID}", handler.HandleUpdateGoal)
		r.Delete("/{goalID}", handler.HandleDeleteGoal)
		r.Post("/{goalID}/tasks", handler.HandleCreateTask)
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/assigned", handler.HandleGetAssignedTasks)
		r.Put("/{taskID}", handler.HandleUpdateTask)
		r.Delete("/{taskID}", handler.HandleDeleteTask)
		r.Put("/{taskID}/assign", handler.HandleAssignTask)
	})
}
