package evap

import "sync"

type Cooling struct {
	mu      sync.Mutex
	running bool
	hours   float64
}

func NewCooling() *Cooling {
	return &Cooling{}
}

func (c *Cooling) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.running = true
}

func (c *Cooling) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.running = false
}

func (c *Cooling) Running() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

func (c *Cooling) AddHours(hours float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hours += hours
}

func (c *Cooling) Hours() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hours
}
