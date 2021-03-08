// Package configs ...
package configs

import (
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Store *Configs

func init() {
	logFile, err := os.OpenFile("./logs/log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	gin.DefaultErrorWriter = io.MultiWriter(logFile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(logFile)
}

func Init() {
	Store = &Configs{}
	v := viper.New()
	v.SetConfigFile("./config.json")
	v.SetConfigType("json")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = v.Unmarshal(Store)
	if err != nil {
		log.Fatalln(err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		err = v.Unmarshal(Store)
		if err != nil {
			log.Fatalln(err)
		}
	})

	validate := validator.New()

	err = validate.Struct(Store)
	if err != nil {
		log.Fatalln(err)
	}
}
