package main

import (
	"context"
	"log"

	"github.com/MacbotX/simplebank_v1/api"
	db "github.com/MacbotX/simplebank_v1/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

var (
	conn *pgxpool.Pool
)

func main() {
	var err error
	conn, err = pgxpool.New(context.Background(), dbSource)
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	defer conn.Close()
}
