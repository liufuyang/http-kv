package cache

type counter struct {
	ch    chan int
	count int
}

func newCounter() *counter {
	ch := make(chan int, 1000000)
	c := counter{ch: ch, count: 0}
	go func() {
		for {
			v := <-c.ch
			c.count += v
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
