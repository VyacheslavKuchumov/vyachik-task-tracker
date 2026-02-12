package api

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/service/tracker"
	"VyacheslavKuchumov/test-backend/service/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	trackerStore := tracker.NewStore(s.db)
	trackerHandler := tracker.NewHandler(trackerStore)

	r.Get("/", trackerHandler.HandleHome)
	r.Get("/login", trackerHandler.HandleLoginPage)
	r.Get("/register", trackerHandler.HandleRegisterPage)
	r.Get("/goals", auth.WithJWTPageAuth(trackerHandler.HandleGoalsPage, userStore))
	r.Get("/tasks", auth.WithJWTPageAuth(trackerHandler.HandleTasksPage, userStore))
	r.Post("/auth/login", userHandler.HandleWebLogin)
	r.Post("/auth/register", userHandler.HandleWebRegister)
	r.Post("/auth/logout", userHandler.HandleLogout)

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/login", userHandler.HandleLogin)
		r.Post("/register", userHandler.HandleRegister)

		r.Route("/goals", func(r chi.Router) {
			r.Get("/", auth.WithJWTAuth(trackerHandler.HandleGetGoals, userStore))
			r.Post("/", auth.WithJWTAuth(trackerHandler.HandleCreateGoal, userStore))
			r.Post("/{goalID}/tasks", auth.WithJWTAuth(trackerHandler.HandleCreateTask, userStore))
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/assigned", auth.WithJWTAuth(trackerHandler.HandleGetAssignedTasks, userStore))
			r.Put("/{taskID}/assign", auth.WithJWTAuth(trackerHandler.HandleAssignTask, userStore))
		})
	})

	r.Route("/htmx", func(r chi.Router) {
		r.Get("/goals", auth.WithJWTAuth(trackerHandler.HandleHTMXGoals, userStore))
		r.Get("/goals/card", auth.WithJWTAuth(trackerHandler.HandleHTMXGoalCard, userStore))
		r.Get("/goals/card/{goalID}", auth.WithJWTAuth(trackerHandler.HandleHTMXGoalCard, userStore))
		r.Post("/goals/save", auth.WithJWTAuth(trackerHandler.HandleHTMXGoalSave, userStore))
		r.Get("/tasks", auth.WithJWTAuth(trackerHandler.HandleHTMXTasks, userStore))
		r.Get("/tasks/card", auth.WithJWTAuth(trackerHandler.HandleHTMXTaskCard, userStore))
		r.Get("/tasks/card/{taskID}", auth.WithJWTAuth(trackerHandler.HandleHTMXTaskCard, userStore))
		r.Post("/tasks/save", auth.WithJWTAuth(trackerHandler.HandleHTMXTaskSave, userStore))
	})

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, r)
}
