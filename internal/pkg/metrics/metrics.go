// Package metrics ..
package metrics

import (
	"app/internal/pkg/sessions"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/IncSW/geoip2"
	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
)

var geoIPReader *geoip2.CityReader
var metrics map[int]map[string]map[string]int
var maxLatency map[string]time.Duration
var metricsMutex *sync.Mutex

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Get(c)
		start := time.Now()

		c.Next()

		path := c.Request.URL.Path

		ua := ua.Parse(c.Request.UserAgent())

		IncrValue("handlers", path)
		IncrValue("os", ua.OS)
		IncrValue("browsers", ua.Name)
		IncrValue("genders", session.Gender)

		latency := time.Since(start)

		if latency > GetMaxLatency(path) {
			SetMaxLatency(path, latency)
		}
		var ip string
		if c.ClientIP() == "::1" {
			ip = "79.141.162.81"
		} else {
			ip = "51.83.70.23"
		}
		record, err := geoIPReader.Lookup(net.ParseIP(ip))
		if err == nil {
			IncrValue("countries", record.Country.Names["en"])
		}
	}
}

func InfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "metrics.html", gin.H{
			"Time":       fmt.Sprintf("From %v:00 to %v:00", time.Now().Hour(), time.Now().Hour()+1),
			"Endpoints":  metrics[time.Now().Hour()]["handlers"],
			"MaxLatency": maxLatency,
			"TotalViews": total(),
			"Browsers":   metrics[time.Now().Hour()]["browsers"],
			"Countries":  metrics[time.Now().Hour()]["countries"],
			"Genders":    metrics[time.Now().Hour()]["genders"],
			"Os":         metrics[time.Now().Hour()]["os"],
		})
	}
}

func total() (total int) {
	metricsMutex.Lock()
	for _, value := range metrics[time.Now().Hour()]["handlers"] {
		total += value
	}
	metricsMutex.Unlock()
	return total
}

func init() {
	var err error
	geoIPReader, err = geoip2.NewCityReaderFromFile("/usr/share/GeoIP/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatalln(err)
	}
	metrics = map[int]map[string]map[string]int{}
	maxLatency = map[string]time.Duration{}
	metricsMutex = &sync.Mutex{}
}

func IncrValue(metricType string, key string) {
	metricsMutex.Lock()
	if metrics[time.Now().Hour()] == nil {
		metrics[time.Now().Hour()] = make(map[string]map[string]int)
	}
	if metrics[time.Now().Hour()][metricType] == nil {
		metrics[time.Now().Hour()][metricType] = make(map[string]int)
	}
	metrics[time.Now().Hour()][metricType][key]++
	metricsMutex.Unlock()
}

func Set(metricType string, key string, value string) {
	metricsMutex.Lock()
	if metrics[time.Now().Hour()] == nil {
		metrics[time.Now().Hour()] = make(map[string]map[string]int)
	}
	if metrics[time.Now().Hour()][metricType] == nil {
		metrics[time.Now().Hour()][metricType] = make(map[string]int)
	}
	metrics[time.Now().Hour()][metricType][key]++
	metricsMutex.Unlock()
}

func GetMaxLatency(path string) (duration time.Duration) {
	metricsMutex.Lock()
	duration = maxLatency[path]
	metricsMutex.Unlock()
	return duration
}

func SetMaxLatency(path string, duration time.Duration) {
	metricsMutex.Lock()
	maxLatency[path] = duration
	metricsMutex.Unlock()
}
