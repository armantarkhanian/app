// Package db ...
package db

import (
	"app/internal/pkg/configs"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func Init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", configs.Store.MySQL.User, configs.Store.MySQL.Password, "%s", configs.Store.MySQL.Database)
	masterDSN := fmt.Sprintf(dsn, configs.Store.MySQL.Master)
	DB, err = sql.Open("mysql", masterDSN)
	if err != nil {
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
}
