package server

import (
	"VyacheslavKuchumov/test-backend/service/auth"
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
	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, s.router())
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)

	trackerStore := tracker.NewStore(s.db)
	trackerHandler := tracker.NewHandler(trackerStore)
	authMiddleware := auth.JWTAuthMiddleware(userStore)
	apiAuthMiddleware := auth.JWTAuthMiddlewareWithExclusions(
		userStore,
		"/api/v1/login",
		"/api/v1/register",
	)

	r.With(authMiddleware).Handle("/swagger/*", httpSwagger.Handler())

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(apiAuthMiddleware)
		user.RegisterRoutes(api, userHandler)
		tracker.RegisterRoutes(api, trackerHandler)
	})

	return r
}
