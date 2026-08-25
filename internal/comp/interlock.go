package comp

import "coldstore/internal/ammonia"

type Interlock struct {
	monitor *ammonia.Monitor
}

func NewInterlock(monitor *ammonia.Monitor) *Interlock {
	return &Interlock{monitor: monitor}
}

func (i *Interlock) Trip(reading float64) bool {
	return reading >= i.monitor.Threshold()
}
