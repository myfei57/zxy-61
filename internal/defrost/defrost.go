package defrost

import (
	"sync"

	"coldstore/internal/fan"
)

type Controller struct {
	mu     sync.Mutex
	fan    *fan.StopUnit
	heater bool
	steps  []string
}

func NewController(fanUnit *fan.StopUnit) *Controller {
	return &Controller{fan: fanUnit, steps: []string{}}
}

func (c *Controller) HeaterOn() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.heater
}

func (c *Controller) Steps() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]string, len(c.steps))
	copy(result, c.steps)
	return result
}
