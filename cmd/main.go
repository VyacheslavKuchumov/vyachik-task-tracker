package main

import (
	"VyacheslavKuchumov/test-backend/cmd/server"
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/db"
	"log"
)

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
