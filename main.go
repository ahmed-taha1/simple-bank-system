package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://taha:root@localhost:5431/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8090"
)

func main() {
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