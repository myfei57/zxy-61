package quota

import "sync"

type Cycle struct {
	mu       sync.Mutex
	limit    int
	consumed int
}

func NewCycle(limit int) *Cycle {
	return &Cycle{limit: limit}
}

func (c *Cycle) Consume(amount int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.consumed+amount > c.limit {
		return false
	}
	c.consumed += amount
	return true
}

func (c *Cycle) Remaining() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.limit - c.consumed
}

func (c *Cycle) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.consumed = 0
}
