package cache

import (
	"sync"
)

type counter struct {
	ch    chan int
	count int
	mux   sync.Mutex
}

func newCounter() *counter {
	ch := make(chan int, 1000000)
	c := counter{ch: ch, count: 0}
	go func() {
		for {
			v := <-c.ch
			c.mux.Lock()
			c.count += v
			c.mux.Unlock()
		}
	}()
	return &c
}

func (c *counter) inc() {
	c.ch <- 1
}

func (c *counter) dec() {
	c.ch <- -1
}

func (c *counter) countSafe() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.count
}

func (c *counter) setSafe(newCount int) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.count = newCount
}
