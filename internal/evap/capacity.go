package evap

import "sync"

type Capacity struct {
	mu       sync.Mutex
	tons     float64
	setpoint float64
}

func NewCapacity(tons float64, setpoint float64) *Capacity {
	return &Capacity{tons: tons, setpoint: setpoint}
}

func (c *Capacity) Update(tons float64, setpoint float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.tons = tons
	c.setpoint = setpoint
}

func (c *Capacity) CoolingCapacity() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.tons * 3.517
}

func (c *Capacity) Setpoint() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setpoint
}

func (c *Capacity) Margin(current float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.setpoint - current
}
