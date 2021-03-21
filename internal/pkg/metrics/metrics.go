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
	BackendIP string                    `json:"backendIP"`
	Year      int                       `json:"year"`
	Month     int                       `json:"month"`
	Day       int                       `json:"day"`
	Hour      int                       `json:"hour"`
	Status    string                    `json:"status"`
	Metrics   map[string]map[string]int `json:"metrics"`
}

func (store *metricStore) Incr(metricType, key string) {
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

func (metricStruct *metricStore) Get(metricType, key string) (value int) {
	metricStruct.RLock()
	value = metricStruct.Metrics[metricType][key]
	metricStruct.RUnlock()
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
	stores[lastHourKey].Status = "completed"
	if err := stores[lastHourKey].saveFile(); err != nil {
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
		log.Println(err)
	}
	bytes, err := ioutil.ReadFile(store.filepath())
	if err != nil {
		log.Println("[ERROR]", err)
	}
	if err := json.Unmarshal(bytes, store); err != nil {
		log.Println("[ERROR]", err)
	}
	if store.Status != "active" && store.Status != "completed" {
		store.Status = "active"
	}
}

func (store *metricStore) saveFile() error {
	bytes, err := json.MarshalIndent(store, "", "\t")
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
		BackendIP: global.BackendIP,
		Year:      year,
		Month:     month,
		Day:       day,
		Hour:      hour,
		Status:    "active",
	}

	stores[lastHourKey].readFromFile()
}

func task() {
	currentStore().saveFile()
}

func Incr(metricType, key string) {
	currentStore().Incr(metricType, key)
}

func Get(metricType, key string) (value int) {
	return currentStore().Get(metricType, key)
}
