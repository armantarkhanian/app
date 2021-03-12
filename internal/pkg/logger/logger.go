// Package logger ...
package logger

import (
	"app/internal/pkg/global"
	"log"
	"os"

	"github.com/rs/zerolog"
)

func Init() {
	logFile, err := os.OpenFile("./logs/main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FATAL] [logger] error opening file:", err)
	}
	accessLog, err := os.OpenFile("./logs/http.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FATAL] [logger] error opening file:", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(logFile)

	global.AccessLog = zerolog.New(accessLog)
}
