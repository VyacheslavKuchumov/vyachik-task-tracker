package db

import (
	"VyacheslavKuchumov/test-backend/config"
	"database/sql"
	"log"
	"net/url"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresStorage(cfg config.Config) (*sql.DB, error) {
	dsn := (&url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   cfg.DBHost + ":" + cfg.DBPort,
		Path:   cfg.DBName,
	}).String() + "?sslmode=" + cfg.DBSSLMode

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open PostgreSQL connection:", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
