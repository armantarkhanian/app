// Package logger ...
package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"app/internal/pkg/about"

	"github.com/rs/zerolog"
)

var (
	mainLog   *os.File
	AccessLog zerolog.Logger
)

func init() {
	var err error
	mainLog, err = os.OpenFile(about.LogDirectory+"main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

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

func Info(v ...interface{}) {
	mainLog.WriteString(output("INFO   ", fmt.Sprint(v...)))
}
func Error(v ...interface{}) {
	mainLog.WriteString(output("ERROR  ", fmt.Sprint(v...)))
}
func Warning(v ...interface{}) {
	mainLog.WriteString(output("WARNING", fmt.Sprint(v...)))
}
func Fatal(v ...interface{}) {
	mainLog.WriteString(output("FATAL  ", fmt.Sprint(v...)))
	os.Exit(1)
}

func Infof(format string, v ...interface{}) {
	mainLog.WriteString(output("INFO   ", fmt.Sprintf(format, v...)))
}
func Errorf(format string, v ...interface{}) {
	mainLog.WriteString(output("ERROR  ", fmt.Sprintf(format, v...)))
}
func Warningf(format string, v ...interface{}) {
	mainLog.WriteString(output("WARNING", fmt.Sprintf(format, v...)))
}
func Fatalf(format string, v ...interface{}) {
	mainLog.WriteString(output("FATAL  ", fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func output(prefix string, data string) string {
	/*
		switch prefix {
		case "FTL":
			data = "FATAL: " + data
		case "ERR":
			data = "ERROR: " + data
		case "   ":
			data = "INFO: " + data
		case "WRN":
			data = "WARNING: " + data
		}*/
	var (
		file string
		line int
		ok   bool
	)
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	caller := file + ":" + strconv.Itoa(line) + ":"
	value := 25 - len(caller)
	if value >= 0 {
		caller = strings.Repeat(".", value) + caller
	} else {
		value = value * (-1)
		caller = caller[value:]
	}
	t := time.Now().UTC().Format("2006/02/01 15:04:05") + " "
	return t + caller + " " + prefix + " " + data + "\n"
}
