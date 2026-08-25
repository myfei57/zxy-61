package service

import (
	"coldstore/internal/room"
)

func (s *Service) Recover() error {
	if err := s.Load(); err != nil {
		return err
	}
	if namespace, ok := s.namespaces.Get("main"); ok {
		s.storeRegistry.Put("namespace", namespace)
	}
	s.storeRegistry.Put("mapping", s.mapping)
	return nil
}

func (s *Service) RoomCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.rooms)
}

func (s *Service) ZoneOf(roomID string) (string, bool) {
	return s.membership.ZoneOf(roomID)
}

func (s *Service) RoomState(roomID string) room.State {
	s.mu.Lock()
	defer s.mu.Unlock()
	if controller, ok := s.roomStates[roomID]; ok {
		return controller.State()
	}
	return room.StateNormal
}

func (s *Service) SetRoomState(roomID string, state room.State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if controller, ok := s.roomStates[roomID]; ok {
		controller.SetState(state)
	}
}

func (s *Service) CoolingCapacity() float64 {
	return s.capacity.CoolingCapacity()
}

func (s *Service) CoolingMargin(current float64) float64 {
	return s.capacity.Margin(current)
}

func (s *Service) FanDuty() (float64, float64, float64) {
	return s.fanDuty.Active(), s.fanDuty.Total(), s.fanDuty.Ratio()
}

func (s *Service) FanRuntime() float64 {
	return s.fanRuntime.Hours()
}

func (s *Service) AlarmState() string {
	return string(s.ammoniaMonitor.State())
}

func (s *Service) AmmoniaPeak() float64 {
	return s.ammoniaLevels.Peak()
}

func (s *Service) AmmoniaAverage() float64 {
	return s.ammoniaLevels.Average()
}

func (s *Service) QuotaWindowSum() int {
	return s.quotaWindow.Sum()
}

func (s *Service) QuotaCycleRemaining() int {
	return s.quotaCycle.Remaining()
}

func (s *Service) ScheduleBand(load float64) int {
	return s.schedule.Band(load)
}

func (s *Service) BandCurrent() int {
	return s.band.Current()
}

func (s *Service) SeriesSlope() float64 {
	return s.series.Slope()
}

func (s *Service) DerivativeRate() float64 {
	return s.derivative.Value()
}

func (s *Service) ValveMode() string {
	return string(s.actuator.Mode())
}
