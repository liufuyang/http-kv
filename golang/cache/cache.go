package cache

import "golang.org/x/sync/syncmap"

type Cache interface {
	Get(key string) string
	Set(key string, value string)
}

type SyncmapCache struct {
	m syncmap.Map
}

func NewSyncmapCache() SyncmapCache {
	return SyncmapCache{m: syncmap.Map{}}
}

func (sc *SyncmapCache) Get(key string) string {
	v, _ := sc.m.Load(key)
	vStr, _ := v.(string)
	return vStr
}

func (sc *SyncmapCache) Set(key string, value string) {
	sc.m.Store(key, value)
}
