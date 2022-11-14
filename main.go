package main

import (
	"database/sql"
	"go-backend-demo/api"
	db "go-backend-demo/db/sqlc"
	"go-backend-demo/util"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	conn, err := sql.Open(config.Db.Driver, config.Db.Source)
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.Server.Addr); err != nil {
		log.Fatal("can not start server:", err)
	}
}
