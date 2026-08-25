package ammonia

func (m *Monitor) Calibrate(threshold float64) {
	m.threshold = threshold
}

func (m *Monitor) Threshold() float64 {
	return m.threshold
}
