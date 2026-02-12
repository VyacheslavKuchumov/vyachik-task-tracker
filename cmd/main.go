package main

import (
	"VyacheslavKuchumov/test-backend/cmd/api"
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/db"
	"log"
)

func main() {
	db, err := db.NewPostgresStorage(config.Envs)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(config.Envs.Port, db)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}

}
