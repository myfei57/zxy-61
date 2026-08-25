package ammonia

import "sync"

type Monitor struct {
	mu        sync.Mutex
	threshold float64
	reading   float64
}

func NewMonitor(threshold float64) *Monitor {
	return &Monitor{threshold: threshold}
}

func (m *Monitor) Read(concentration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reading = concentration
}

func (m *Monitor) Reading() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.reading
}
