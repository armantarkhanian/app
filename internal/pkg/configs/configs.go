// Package configs ...
package configs

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/go-playground/validator/v10"
)

var Store *Configs

func Init() {
	bytes, err := ioutil.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(bytes, &Store)
	if err != nil {
		log.Fatalln(err)
	}

	validate := validator.New()

	err = validate.Struct(Store)
	if err != nil {
		log.Fatalln(err)
	}
}
