package service

func (s *Service) CheckInterlock(reading float64) bool {
	s.ammoniaMonitor.Read(reading)
	tripped := s.interlock.Trip(reading)
	s.auditLog.Add("comp", "interlock-check", "ammonia")
	return tripped
}

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
