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

// ReadingAndThreshold returns the current reading and threshold as a single
// atomic snapshot, so callers that need both values (e.g. for a consistent
// status/log line) cannot observe a reading from one moment and a threshold
// from another.
func (m *Monitor) ReadingAndThreshold() (float64, float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.reading, m.threshold
}

func (m *Monitor) Reading() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.reading
}
