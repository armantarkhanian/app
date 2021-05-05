// Package logger ...
package logger

import (
	"log"
	"os"
	"strconv"
	"time"

	"app/internal/pkg/about"

	"github.com/rs/zerolog"
)

var (
	AccessLog zerolog.Logger
)

func init() {
	mainLog, err := os.OpenFile(about.LogDirectory+"main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(mainLog)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)

	accessLog, err := os.OpenFile(about.LogDirectory+"access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	zerolog.MessageFieldName = "msg"
	zerolog.TimestampFieldName = "t"
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.LevelFieldName = "l"
	zerolog.TimestampFunc = time.Now().UTC
	zerolog.CallerFieldName = "c"
	zerolog.CallerSkipFrameCount = 3
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	AccessLog = zerolog.New(accessLog).With().Timestamp().Logger().With().Caller().Logger().With().Str("host", "192.143.34.54").Logger()
}
