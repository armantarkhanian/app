// Package about ...
package about

import (
	"log"
	"os"
)

const (
	LogDirectory    string = "/var/log/app/"
	ConfigDirectory string = "/etc/app/"
)

func init() {
	if err := os.MkdirAll(LogDirectory, 0777); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll(ConfigDirectory, 0777); err != nil {
		log.Fatalln(err)
	}
}
