package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/natanaelrusli/go-bank/util"
)

var testQueries *Queries
var testMockQueries *Queries
var mocksql sqlmock.Sqlmock
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("failed reading config file", err)
	}

	_, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("failed creating new mock sql db")
	}

	defer mockDB.Close()

	testMockQueries = New(mockDB)
	mocksql = mock

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("failed connecting to testDB")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
