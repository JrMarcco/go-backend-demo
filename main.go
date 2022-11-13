package main

import (
	"database/sql"
	"go-backend-demo/api"
	db "go-backend-demo/db/sqlc"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	conn, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/simple_bank?parseTime=true")
	if err != nil {
		log.Fatal("can not connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(":8080"); err != nil {
		log.Fatal("can not start server:", err)
	}
}
