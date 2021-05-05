// Package geoip ...
package geoip

import (
	"log"
	"net"

	"github.com/IncSW/geoip2"
)

func CountryCodeByIP(ip string) string {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return ""
	}
	record, err := reader.Lookup(netIP)
	if err != nil {
		if err != geoip2.ErrNotFound {
			log.Println("[ERROR] [geoip]", err)
		}
		return ""
	}
	return record.Country.ISOCode
}
