package server

import (
	"VyacheslavKuchumov/test-backend/service/tracker"
	"VyacheslavKuchumov/test-backend/service/user"
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/swagger/*", httpSwagger.Handler())

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	trackerStore := tracker.NewStore(s.db)
	trackerHandler := tracker.NewHandler(trackerStore)

	r.Route("/api/v1", func(api chi.Router) {
		user.RegisterRoutes(api, userHandler)
		tracker.RegisterRoutes(api, trackerHandler, userStore)
	})

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, r)
}
