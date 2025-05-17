package main

import (
	"context"
	"log"

	"github.com/MacbotX/simplebank_v1/api"
	db "github.com/MacbotX/simplebank_v1/db/sqlc"
	"github.com/MacbotX/simplebank_v1/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	conn *pgxpool.Pool
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err = pgxpool.New(context.Background(), config.DBSource)
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	defer conn.Close()
}
