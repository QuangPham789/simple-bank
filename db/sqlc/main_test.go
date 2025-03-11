package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	// dbTest, err = pgx.Connect(context.Background(), dbSource)
	testDB, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	conn, err := testDB.Acquire(context.Background()) // Acquire a new connection from the pool
	if err != nil {
		panic(err)
	}
	defer conn.Release()
	if err != nil {
		log.Fatal("can't connect to database:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
