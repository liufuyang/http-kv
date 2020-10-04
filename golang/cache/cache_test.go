package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	shortCache := NewSyncmapCache(time.Duration(1) * time.Second)

	value1 := shortCache.Get("key1")
	size := shortCache.Size()
	if value1 != "" {
		t.Errorf("value1 not empty, having value '%v'", value1)
	}
	if size != 0 {
		t.Errorf("size is not 0 but '%v'", size)
	}

	shortCache.Set("key1", "value1")
	value1 = shortCache.Get("key1")
	time.Sleep(time.Duration(100) * time.Millisecond)
	size = shortCache.Size()
	if value1 != "value1" {
		t.Errorf("value1 is not 'value1' but '%v'", value1)
	}
	if size != 1 {
		t.Errorf("size is not 1 but '%v'", size)
	}

	time.Sleep(time.Duration(2) * time.Second)
	value1 = shortCache.Get("key1")
	size = shortCache.Size()
	if value1 != "" {
		t.Errorf("value1 not empty, having value '%v'", value1)
	}
	if size != 0 {
		t.Errorf("size is not 0 but '%v'", size)
	}
}

func TestCounter(t *testing.T) {
	c := newCounter()

	c.inc()
	c.inc()

	time.Sleep(time.Duration(100) * time.Millisecond)
	if c.count != 2 {
		t.Errorf("counter count is not 2 but '%v'", c.count)
	}
}
