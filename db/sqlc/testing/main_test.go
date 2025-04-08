package db

import (
	"database/sql"
	"log"
	"os"
	db "simplebank/db/sqlc"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://taha:root@localhost:5431/simple_bank?sslmode=disable"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	testQueries = db.New(conn)

	os.Exit(m.Run())
}
