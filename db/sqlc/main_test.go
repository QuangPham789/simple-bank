package db

import (
	"context"
	"github.com/QuangPham789/simple-bank/util"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	// dbTest, err = pgx.Connect(context.Background(), dbSource)
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
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
