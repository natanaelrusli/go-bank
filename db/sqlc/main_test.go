package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testMockQueries *Queries
var mocksql sqlmock.Sqlmock
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	_, err = sql.Open(dbDriver, dbSource)
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

	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("failed connecting to testDB")
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
