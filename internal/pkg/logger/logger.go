// Package logger ...
package logger

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var (
	Zerolog zerolog.Logger
)

func Init() {
	logFile, err := os.OpenFile("./logs/main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FATAL] error opening file:", err)
	}
	accessLog, err := os.OpenFile("./logs/http.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("[FATAL] error opening file:", err)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(logFile)

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// gin.DefaultWriter = accessLog
	Zerolog = zerolog.New(accessLog)
	gin.DefaultErrorWriter = logFile
}
