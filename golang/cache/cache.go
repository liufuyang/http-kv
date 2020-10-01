package cache

import (
	"fmt"
	"time"

	"golang.org/x/sync/syncmap"
)

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
	expireDuration time.Duration
	m              syncmap.Map
}

func NewSyncmapCache(expireDuration time.Duration) *SyncmapCache {
	cache := SyncmapCache{m: syncmap.Map{}, expireDuration: expireDuration}
	cache.vaccum()
	return &cache
}

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

func (sc *SyncmapCache) Set(key string, value string) {
	v := Value{value: value, timestamp: time.Now()}
	sc.m.Store(key, v)
}

func (sc *SyncmapCache) Size() int {
	length := 0
	sc.m.Range(func(key, _ interface{}) bool {
		length++
		return true
	})
	return length
}

func (sc *SyncmapCache) vaccum() {
	go func() {
		for {
			ms := sc.expireDuration.Milliseconds()
			size := sc.Size()
			var sleepMs int64
			if size <= 1 {
				sleepMs = ms
			} else {
				sleepMs = ms / (int64)(size)
			}

			// Debug print
			fmt.Println("size: ", size)
			fmt.Println("sleepMs: ", sleepMs)

			if size == 0 {
				time.Sleep(time.Duration(sleepMs) * time.Millisecond)
			} else {
				sc.m.Range(func(key, v interface{}) bool {
					time.Sleep(time.Duration(sleepMs) * time.Millisecond)

					value := v.(Value)
					if time.Now().After(value.timestamp.Add(sc.expireDuration)) {
						// Debug print
						fmt.Println("deleting key:", key)
						sc.m.Delete(key)
					}
					return true
				})
			}
		}
	}()
}
