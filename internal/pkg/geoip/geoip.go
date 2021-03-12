// Package geoip ...
package geoip

import (
	"log"
	"net"

	"github.com/IncSW/geoip2"
)

var (
	reader *geoip2.CityReader
)

func init() {
	var err error
	reader, err = geoip2.NewCityReaderFromFile("/usr/share/GeoIP/GeoLite2-City.mmdb")
	if err != nil {
		log.Println("[FATAL] [geoip]", err)
	}
}

func CountryCodeByIP(ip string) string {
	record, err := reader.Lookup(net.ParseIP(ip))
	if err != nil {
		if err != geoip2.ErrNotFound {
			log.Println("[ERROR] [geoip]", err)
		}
		return ""
	}
	return record.Country.ISOCode
}
