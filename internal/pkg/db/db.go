// Package db ...
package db

import (
	"app/internal/pkg/configs"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type dbStore struct {
	sync.Mutex
	master            *sql.DB
	slaves            []*sql.DB
	roundRobinCounter int
}

var (
	db dbStore
)

func Master() *sql.DB {
	return db.master
}

func Slave() *sql.DB {
	if len(db.slaves) == 0 {
		return db.master
	}

	db.Lock()
	r := db.slaves[db.roundRobinCounter]
	db.roundRobinCounter = (db.roundRobinCounter + 1) % len(db.slaves)
	db.Unlock()

	return r
}

func Init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", configs.Store.MySQL.User, configs.Store.MySQL.Password, "%s", configs.Store.MySQL.Database)
	masterDSN := fmt.Sprintf(dsn, configs.Store.MySQL.Master)
	db.master, err = sql.Open("mysql", masterDSN)
	if err != nil {
		panic(err)
	}
	if err = db.master.Ping(); err != nil {
		panic(err)
	}
	db.master.SetConnMaxLifetime(time.Minute * 3)
	db.master.SetMaxOpenConns(10)
	db.master.SetMaxIdleConns(10)
	for _, slaveHost := range configs.Store.MySQL.Slaves {
		slaveDSN := fmt.Sprintf(dsn, slaveHost)
		slaveDB, err := sql.Open("mysql", slaveDSN)
		if err != nil {
			panic(err)
		}
		if err = slaveDB.Ping(); err != nil {
			panic(err)
		}
		slaveDB.SetConnMaxLifetime(time.Minute * 3)
		slaveDB.SetMaxOpenConns(10)
		slaveDB.SetMaxIdleConns(10)
		db.slaves = append(db.slaves, slaveDB)
	}
}
