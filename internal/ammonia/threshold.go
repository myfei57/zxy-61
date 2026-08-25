package ammonia

// Calibrate publishes a new ammonia threshold. It must take the monitor lock so
// that concurrent interlock reads (Threshold/Evaluate/ReadingAndThreshold)
// observe the updated value with a proper happens-before relationship instead
// of a stale pre-calibration threshold that could let an over-limit
// concentration miss the trip.
func (m *Monitor) Calibrate(threshold float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.threshold = threshold
}

// Threshold returns the currently configured ammonia threshold. It shares the
// monitor lock with Calibrate so the interlock never reads the threshold field
// concurrently with a calibration write.
func (m *Monitor) Threshold() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.threshold
}
