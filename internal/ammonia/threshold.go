package ammonia

func (m *Monitor) Calibrate(threshold float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.threshold = threshold
}

func (m *Monitor) Threshold() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.threshold
}
