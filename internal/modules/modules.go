// Package modules ...
package modules

import (
	"app/internal/pkg/geoip"
	"log"
)

type Modules struct {
	RedisCluster  bool
	InnoDBCluster bool
	GeoIP         bool
}

var (
	m = Modules{}
)

func Init(m Modules) {
	if m.GeoIP {
		if err := geoip.InitModule(); err != nil {
			log.Fatalln(err)
		}
	}
}
