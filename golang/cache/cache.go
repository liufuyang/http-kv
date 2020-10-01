package cache

import (
	"time"

	"golang.org/x/sync/syncmap"
)

const ExpireDurationStr = "30m"

type Value struct {
	value     string
	timestamp time.Time
}

type Cache interface {
	Get(key string) string
	Set(key string, value string)
	Size() int
}

type SyncmapCache struct {
	m syncmap.Map
}

func NewSyncmapCache() SyncmapCache {
	return SyncmapCache{m: syncmap.Map{}}
}

func (sc *SyncmapCache) Get(key string) string {
	v, ok := sc.m.Load(key)
	if !ok {
		return ""
	}
	value := v.(Value)

	expireDuration, _ := time.ParseDuration(ExpireDurationStr)
	if time.Now().After(value.timestamp.Add(expireDuration)) {
		sc.m.Delete(key)
		return ""
	}
	return value.value
}

func (sc *SyncmapCache) Set(key string, value string) {
	v := Value{value: value, timestamp: time.Now()}
	sc.m.Store(key, v)
}

func (sc *SyncmapCache) Size() int {
	length := 0
	sc.m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}