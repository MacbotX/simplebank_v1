package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"

var (
	testQueries *Queries
 	conn *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var err error
	conn, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	
	defer conn.Close()
	testQueries = New(conn)
	// Run the tests and exit appropriately
	os.Exit(m.Run())
}
