package comp

import "coldstore/internal/ammonia"

// Interlock evaluates ammonia concentration against the configured threshold to
// decide whether the safety trip must engage.
type Interlock struct {
	monitor *ammonia.Monitor
}

func NewInterlock(monitor *ammonia.Monitor) *Interlock {
	return &Interlock{monitor: monitor}
}

// Trip publishes the supplied concentration reading to the monitor and, within
// the same critical section that records it, evaluates whether it trips the
// interlock against the threshold currently in effect. Doing both under one
// lock guarantees the recorded reading and the threshold used for the trip
// decision are a consistent pair: a concurrent Calibrate cannot publish a new
// threshold in between, so the interlock never judges a post-calibration
// reading against a stale, pre-calibration threshold.
func (i *Interlock) Trip(reading float64) bool {
	return i.monitor.Evaluate(reading)
}
