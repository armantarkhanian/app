// Package configs ...
package configs

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var Store *Configs

func Init() {
	Store = &Configs{}
	v := viper.New()
	v.SetConfigFile("./config.json")
	v.SetConfigType("json")

	err := v.ReadInConfig()
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}

	err = v.Unmarshal(Store)
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		err = v.Unmarshal(Store)
		if err != nil {
			log.Fatalln("[FATAL]", err)
		}
		log.Println("[INFO]", "<config.json> was successfully re-read")
	})

	validate := validator.New()

	err = validate.Struct(Store)
	if err != nil {
		log.Fatalln("[FATAL]", err)
	}
}
