package service

func (s *Service) Defrost() error {
	if err := s.defrostController.Start(); err != nil {
		return err
	}
	s.fanUnit.Stop()
	s.auditLog.Add("console", "defrost", "start")
	return nil
}

func (s *Service) RunDefrost() error {
	if err := s.defrostRun.Run(s.supplyValve); err != nil {
		return err
	}
	s.auditLog.Add("console", "defrost-run", "run")
	return nil
}
