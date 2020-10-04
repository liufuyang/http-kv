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
	c              *counter
}

// NewSyncmapCache is for creating a new SyncmapCache
func NewSyncmapCache(expireDuration time.Duration) *SyncmapCache {
	cache := SyncmapCache{m: syncmap.Map{}, expireDuration: expireDuration, c: newCounter()}
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
		sc.c.dec()
		cacheSizeGauge.Dec()
		return ""
	}
	return value.value
}

// Set is for impl of Cache for SyncmapCache - set value
func (sc *SyncmapCache) Set(key string, value string) {
	v := Value{value: value, timestamp: time.Now()}
	if _, ok := sc.m.Load(key); !ok {
		sc.c.inc()
		cacheSizeGauge.Inc()
	}
	sc.m.Store(key, v)
}

// Size to return the current element size in Cache
func (sc *SyncmapCache) Size() int {
	// TODO remove debug print
	// log.Info("For loop calculate size: ", sc.SizeOld())
	log.Info("Cache size estimation: ", sc.c.count)
	return sc.c.count
}

// SizeOld to return the current element size in Cache, the loop way
// Deprecated
func (sc *SyncmapCache) SizeOld() int {
	length := 0
	sc.m.Range(func(key, _ interface{}) bool {
		length++
		return true
	})
	return length
}

// vaccum method is used for SyncmapCache clean up expired key
// We set up a minimum sleep time as 2ms durring the key loop to improve overall cache performace during high load
func (sc *SyncmapCache) vaccum() {
	log.Info("cache expire time: ", sc.expireDuration)
	vaccumCycleDurationMs := int64(60000)
	minimumSleepMs := int64(2)
	go func() {
		for {
			size := sc.Size()
			log.Debug("cache size: ", size)

			var sleepMs int64
			if size <= 1 {
				sleepMs = vaccumCycleDurationMs
			} else {
				sleepMs = max(vaccumCycleDurationMs/(int64)(size), minimumSleepMs) // Without slowing down (setting a max 10ms sleep time) it can reduce performace with many key expires at the same time
			}

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
						sc.c.dec()
						cacheSizeGauge.Dec()
					}
					return true
				})
			}
		}
	}()
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
