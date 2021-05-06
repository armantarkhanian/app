package main

import (
	"app/internal/pkg/logger"
	"database/sql"
	"errors"
)

func main() {
	logger.Info("server is running at port :8080")
	logger.Infof("%s", "server is running at port :8080")
	_, err := sql.Open("mysql", "faq")
	if err != nil {
		logger.Errorf("ERROR: %v", err)
	}
	logger.Warningf("can't connect to telegram bot server: %v", errors.New("http timeout deadline"))
	logger.Info("server is running at port :8080")
	_, err = sql.Open("mysql", "faq")
	if err != nil {
		logger.Fatal(err)
	}
}
