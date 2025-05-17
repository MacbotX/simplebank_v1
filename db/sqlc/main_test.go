package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/MacbotX/simplebank_v1/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testQueries *Queries
	testDB      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	defer testDB.Close()
	testQueries = New(testDB)
	// Run the tests and exit appropriately
	os.Exit(m.Run())
}
