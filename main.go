package main

import (
	log "app/internal/pkg/logger"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Info("Trying to connect to mysql:")
	db, err := sql.Open("mysql", "msandbox:msandbox@tcp(127.0.0.1:26223)/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Info("Successfull connection to InnoDB Cluster")
}
