package main

import (
	"context"
	"github.com/QuangPham789/simple-bank/api"
	db "github.com/QuangPham789/simple-bank/db/sqlc"
	"github.com/QuangPham789/simple-bank/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)

	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("can't create server: %v", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}
