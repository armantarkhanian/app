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
	v.SetConfigFile("./configs/config.json")
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
