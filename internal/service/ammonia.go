package service

// CheckInterlock records the supplied ammonia concentration reading and
// evaluates the safety interlock against it. Both the reading publication and
// the trip evaluation happen under the monitor's lock as a single atomic
// step (via Interlock.Trip -> Monitor.Evaluate), so a concurrent
// CalibrateAmmonia cannot leave the interlock judging the new reading against
// a stale, pre-calibration threshold.
func (s *Service) CheckInterlock(reading float64) bool {
	tripped := s.interlock.Trip(reading)
	s.auditLog.Add("comp", "interlock-check", "ammonia")
	return tripped
}

// CalibrateAmmonia publishes a new ammonia concentration threshold. The
// threshold is published under the monitor lock so that any in-flight or
// subsequent interlock evaluation observes the new value with a proper
// happens-before relationship.
func (s *Service) CalibrateAmmonia(threshold float64) {
	s.ammoniaMonitor.Calibrate(threshold)
	s.auditLog.Add("console", "ammonia-calibrate", "threshold")
}

func (s *Service) FanVerdict(temperature float64) bool {
	return s.fanVerdict.Judge(temperature)
}

func (s *Service) QuotaAllow(amount int) bool {
	return s.quotaTracker.Reserve(amount)
}

func (s *Service) ResetQuota() {
	s.quotaTracker.Reset()
}
