package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/natanaelrusli/go-bank/api"
	db "github.com/natanaelrusli/go-bank/db/sqlc"
	"github.com/natanaelrusli/go-bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
