package cache

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/syncmap"
)

var (
	cacheSizeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cache_size",

		Help: "The number of HTTP requests GET on / served in the last second",
	})
)

// Value type for the kv cache
type Value struct {
	value     string
	timestamp time.Time
}

// Cache is an interface to allow easily switch to a different impl later on
type Cache interface {
	Get(key string) string
	Set(key string, value string)
	Size() int
}

// SyncmapCache is a simple Cache impl with syncmap.Map
type SyncmapCache struct {
	expireDuration time.Duration
	m              syncmap.Map
}

// NewSyncmapCache is for creating a new SyncmapCache
func NewSyncmapCache(expireDuration time.Duration) *SyncmapCache {
	cache := SyncmapCache{m: syncmap.Map{}, expireDuration: expireDuration}
	cache.vaccum()
	return &cache
}

// Get is for impl of Cache for SyncmapCache - get value
func (sc *SyncmapCache) Get(key string) string {
	v, ok := sc.m.Load(key)
	if !ok {
		return ""
	}
	value := v.(Value)

	if time.Now().After(value.timestamp.Add(sc.expireDuration)) {
		sc.m.Delete(key)
		return ""
	}
	return value.value
}

// Set is for impl of Cache for SyncmapCache - set value
func (sc *SyncmapCache) Set(key string, value string) {
	v := Value{value: value, timestamp: time.Now()}
	sc.m.Store(key, v)
}

// Size to return the current element size in Cache
func (sc *SyncmapCache) Size() int {
	length := 0
	sc.m.Range(func(key, _ interface{}) bool {
		length++
		return true
	})
	return length
}

// vaccum method is used for SyncmapCache clean up expired key
func (sc *SyncmapCache) vaccum() {
	log.Info("cache expire time: ", sc.expireDuration)
	var vaccumCycleDurationMs int64
	vaccumCycleDurationMs = 60000
	go func() {
		for {
			size := sc.Size()
			var sleepMs int64
			if size <= 1 {
				sleepMs = vaccumCycleDurationMs
			} else {
				sleepMs = vaccumCycleDurationMs / (int64)(size)
			}
			log.Debug("cache size: ", size)
			cacheSizeGauge.Set(float64(size))

			if size == 0 {
				log.Debug("vaccum sleep: ", sleepMs, "ms")
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
			} else {
				sc.m.Range(func(key, v interface{}) bool {
					log.Debug("vaccum sleep: ", sleepMs, "ms")
					time.Sleep(time.Duration(sleepMs) * time.Millisecond) // sleep here for some time to reduce loop frequency

					value := v.(Value)
					if time.Now().After(value.timestamp.Add(sc.expireDuration)) {
						log.Debug("deleting key:", key)
						sc.m.Delete(key)
					}
					return true
				})
			}
		}
	}()
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
