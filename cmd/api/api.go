package api

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/service/product"
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
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/login", userHandler.HandleLogin)
		r.Post("/register", userHandler.HandleRegister)

		r.Route("/product", func(r chi.Router) {
			r.Get("/", auth.WithJWTAuth(productHandler.HandleGetProducts, userStore))
			r.Post("/", auth.WithJWTAuth(productHandler.HandleCreateProduct, userStore))
		})
	})

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, r)
}
