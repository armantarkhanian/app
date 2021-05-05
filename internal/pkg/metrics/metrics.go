// Package metrics ..
package metrics

import (
	"app/internal/pkg/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jasonlvhit/gocron"
)

type metricStore struct {
	sync.RWMutex
	BackendID int                       `json:"backendID"`
	Year      int                       `json:"year"`
	Month     int                       `json:"month"`
	Day       int                       `json:"day"`
	Hour      int                       `json:"hour"`
	Metrics   map[string]map[string]int `json:"metrics"`
}

func (store *metricStore) incr(metricType, key string) {
	store.Lock()
	if store.Metrics == nil {
		store.Metrics = make(map[string]map[string]int)
	}
	if store.Metrics[metricType] == nil {
		store.Metrics[metricType] = make(map[string]int)
	}
	store.Metrics[metricType][key]++
	store.Unlock()
}

func (store *metricStore) get(metricType, key string) (value int) {
	store.RLock()
	value = store.Metrics[metricType][key]
	store.RUnlock()
	return value
}

// hourKey is a string in format year/month/day/hour
type hourKey string

var (
	lastHourKey hourKey
	stores      map[hourKey]*metricStore
)

func currentStore() *metricStore {
	if lastHourKey == currentHourKey() {
		return stores[lastHourKey]
	}

	err := stores[lastHourKey].saveFile()

	if err != nil {
		log.Println("[ERROR]", err)
	}

	delete(stores, lastHourKey)

	initCurrentStore()

	return stores[lastHourKey]
}

func currentHourKey() hourKey {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()
	return hourKey(fmt.Sprintf("%v/%v/%v/%v", year, month, day, hour))
}

func (store *metricStore) hourKey() hourKey {
	return hourKey(fmt.Sprintf("%v/%v/%v/%v", store.Year, store.Month, store.Day, store.Hour))
}

func (store *metricStore) filepath() string {
	return "./metrics/" + string(store.hourKey()) + ".json"
}

func (store *metricStore) dir() string {
	array := strings.Split(string(store.hourKey()), "/")
	if len(array) == 4 {
		return strings.Join(array[:3], "/")
	}
	return ""
}

func (store *metricStore) readFromFile() {
	if err := os.MkdirAll("./metrics/"+store.dir(), 0755); err != nil && err != os.ErrExist {
		log.Println("[ERROR]", err)
	}
	bytes, err := ioutil.ReadFile(store.filepath())
	if err != nil {
		log.Println("[ERROR]", err)
	}
	if err := json.Unmarshal(bytes, store); err != nil {
		log.Println("[ERROR]", err)
	}
}

func (store *metricStore) saveFile() error {
	bytes, err := json.Marshal(store)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(store.filepath(), bytes, 0644); err != nil {
		return err
	}
	return nil
}

func init() {
	stores = make(map[hourKey]*metricStore)

	initCurrentStore()

	go func() {
		gocron.Every(30).Seconds().Do(task)
		<-gocron.Start()
	}()
}

func initCurrentStore() {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()
	hour := now.Hour()

	lastHourKey = hourKey(fmt.Sprintf("%v/%v/%v/%v", year, month, day, hour))

	stores[lastHourKey] = &metricStore{
		BackendID: global.BackendID,
		Year:      year,
		Month:     month,
		Day:       day,
		Hour:      hour,
	}

	stores[lastHourKey].readFromFile()
}

func task() {
	if err := currentStore().saveFile(); err != nil {
		log.Println("[ERROR]", err)
	}
}

func Incr(metricType, key string) {
	currentStore().incr(metricType, key)
}

func Get(metricType, key string) (value int) {
	return currentStore().get(metricType, key)
}
