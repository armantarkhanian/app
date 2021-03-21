// Package metrics ..
package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	stores[lastHourKey].saveFile()
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

func (store *metricStore) Filepath() string {
	return "./metrics/" + string(store.hourKey()) + ".json"
}
func (store *metricStore) Dir() string {
	array := strings.Split(string(store.hourKey()), "/")
	if len(array) == 4 {
		return strings.Join(array[:3], "/")
	}
	return ""
}
func (store *metricStore) readFromFile() {
	os.MkdirAll("./metrics/"+store.Dir(), 0755)
	bytes, _ := ioutil.ReadFile(store.Filepath())
	json.Unmarshal(bytes, store)
	if store.Status != "active" && store.Status != "completed" {
		store.Status = "active"
	}
}

func (store *metricStore) saveFile() error {
	bytes, err := json.MarshalIndent(store, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(store.Filepath(), bytes, 0644); err != nil {
		return err
	}
	fmt.Println("File", store.Filepath(), "saved.")
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
		Year:   year,
		Month:  month,
		Day:    day,
		Hour:   hour,
		Status: "active",
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
