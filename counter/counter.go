package counter

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value uint64
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Add(delta uint) (err error) {
	c.mu.Lock()
	defer func() {

		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("%s", pErr)
		}

		c.mu.Unlock()
	}()
	c.value += uint64(delta)
	return nil
}

func (c *Counter) Dec(delta uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.value < uint64(delta) {
		return fmt.Errorf("")
	}
	c.value -= uint64(delta)
	return nil
}

func (c *Counter) Value() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}
