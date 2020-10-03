package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	shortCache := NewSyncmapCache(time.Duration(1) * time.Second)

	value1 := shortCache.Get("key1")
	if value1 != "" {
		t.Errorf("value1 not empty, having value '%v'", value1)
	}

	shortCache.Set("key1", "value1")
	value1 = shortCache.Get("key1")
	if value1 != "value1" {
		t.Errorf("value1 is not 'value1' but '%v'", value1)
	}

	time.Sleep(time.Duration(2) * time.Second)
	value1 = shortCache.Get("key1")
	if value1 != "" {
		t.Errorf("value1 not empty, having value '%v'", value1)
	}
}
