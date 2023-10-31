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

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("failed creating new mock sql db")
	}

	defer mockDB.Close()

	testMockQueries = New(mockDB)
	mocksql = mock

	os.Exit(m.Run())
}
