package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	_ "github.com/lib/pq"
)


func main() {
	// Load configuration from the environment variables or a config file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	dbDriver := config.DBDriver
	dbSource := config.DBSource
	serverAddress := config.ServerAddress

	dbConnection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	store := db.NewStore(dbConnection)
	server := api.CreateNewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}