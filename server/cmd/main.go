package main

import (
	"VyacheslavKuchumov/test-backend/cmd/server"
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/db"
	_ "VyacheslavKuchumov/test-backend/docs"
	"log"
)

// @title Task Tracker API
// @version 1.0
// @description REST API for task tracking with goals and assignments
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db, err := db.NewPostgresStorage(config.Envs)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(config.Envs.Port, db)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
