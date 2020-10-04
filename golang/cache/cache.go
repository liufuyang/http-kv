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
		return ""
	}
	return value.value
}

// Set is for impl of Cache for SyncmapCache - set value
func (sc *SyncmapCache) Set(key string, value string) {
	v := Value{value: value, timestamp: time.Now()}
	if _, ok := sc.m.Load(key); !ok {
		// This should be the only place that counter is increased
		sc.c.inc()
		cacheSizeGauge.Inc()
	}
	sc.m.Store(key, v)
}

// Size to return the approximate current element size in Cache
func (sc *SyncmapCache) Size() int {
	size := sc.c.countSafe()
	return size
}

// SizePrecise to return the current element size in Cache, the loop way
func (sc *SyncmapCache) SizePrecise() int {
	length := 0
	sc.m.Range(func(key, _ interface{}) bool {
		length++
		return true
	})
	return length
}

// vaccum method is used for SyncmapCache clean up expired key
// We set up a minimum sleep time as 1000ns durring the key loop to improve overall cache performace during high load
func (sc *SyncmapCache) vaccum() {
	log.Info("cache expire time: ", sc.expireDuration)
	vaccumCycleSleepMs := min(int64(10000), sc.expireDuration.Milliseconds())
	vaccumKeyLoopSleepNs := int64(1000)
	go func() {
		for {
			sc.m.Range(func(key, v interface{}) bool {
				value := v.(Value)
				if time.Now().After(value.timestamp.Add(sc.expireDuration)) {
					log.Debug("deleting key:", key)
					// This should be the only place counter decreased
					sc.m.Delete(key)
					sc.c.dec()
					cacheSizeGauge.Dec()
					log.Debug("vaccum key-loop sleep: ", vaccumKeyLoopSleepNs, "ns")
					time.Sleep(time.Duration(vaccumKeyLoopSleepNs) * time.Nanosecond) // sleep here for some time to reduce loop frequency
				}
				return true
			})

			// Adjust size for counter and metric gauge
			size := sc.SizePrecise()
			log.Info("Vaccum cycle: Cache size from loop: ", size)
			log.Info("Vaccum cycle: Cache size from counter: ", sc.c.countSafe())
			log.Info("Readjust counter size to: ", size)
			sc.c.setSafe(size)
			cacheSizeGauge.Set(float64(size))

			log.Debug("vaccum sleep: ", vaccumCycleSleepMs, "ms")
			time.Sleep(time.Duration(vaccumCycleSleepMs) * time.Millisecond) // sleep here for some time to reduce loop frequency
		}
	}()
}

func min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
