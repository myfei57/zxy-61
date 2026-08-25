package ammonia

// Evaluate publishes the supplied concentration reading and, within the same
// critical section that records it, decides whether it trips the interlock
// against the threshold currently in effect. Doing both under one lock
// guarantees the recorded reading and the threshold used for the trip decision
// are a consistent pair: a concurrent Calibrate cannot publish a new threshold
// in between, so the interlock never judges a post-calibration reading against
// a stale, pre-calibration threshold.
func (m *Monitor) Evaluate(reading float64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reading = reading
	return reading >= m.threshold
}
