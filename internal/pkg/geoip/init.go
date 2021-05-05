// Package geoip ...
package geoip

import (
	"log"

	"github.com/IncSW/geoip2"
)

var (
	reader       = new(geoip2.CityReader)
	databaseFile = "/usr/share/GeoIP/GeoLite2-City.mmdb"
)

func init() {
	var err error
	reader, err = geoip2.NewCityReaderFromFile(databaseFile)
	if err != nil {
		log.Fatalln(err)
	}
}
