package main

import (
	_ "app/internal/pkg/logger"
	"database/sql"
	"log"
)

func main() {
	log.Println("server is running at port :8080")
	_, err := sql.Open("mysql", "faq")
	if err != nil {
		log.Println(err)
	}
	log.Println("server is running at port :8080")
	_, err = sql.Open("mysql", "faq")
	if err != nil {
		log.Println(err)
		log.Fatalln(err)
	}
}
