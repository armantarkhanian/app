// Package geoip ...
package geoip

import (
	"net"

	"github.com/IncSW/geoip2"
	"github.com/gin-gonic/gin"
)

var (
	reader *geoip2.CityReader
)

func init() {
	var err error
	reader, err = geoip2.NewCityReaderFromFile("/usr/share/GeoIP/GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("countryCode", CountryCodeByIP(c.ClientIP()))
		c.Next()
	}
}

func CountryCodeByIP(ip string) string {
	record, err := reader.Lookup(net.ParseIP(ip))
	if err != nil {
		return ""
	}
	return record.Country.ISOCode
}
