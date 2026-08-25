package service

import "time"

func (s *Service) DetailedStatus() map[string]any {
	fanActive, fanTotal, fanRatio := s.FanDuty()
	return map[string]any{
		"room_count":         s.RoomCount(),
		"series_count":       s.series.Count(),
		"series_average":     s.series.Average(),
		"series_slope":       s.SeriesSlope(),
		"derivative_rate":    s.DerivativeRate(),
		"band_current":       s.BandCurrent(),
		"band_count":         s.ScheduleBandCount(),
		"schedule_highest":   s.ScheduleHighestBand(0),
		"cooling_capacity":   s.CoolingCapacity(),
		"cooling_setpoint":   s.capacity.Setpoint(),
		"cooling_margin":     s.CoolingMargin(0),
		"cooling_running":    s.cooling.Running(),
		"cooling_hours":      s.cooling.Hours(),
		"fan_duty_active":    fanActive,
		"fan_duty_total":     fanTotal,
		"fan_duty_ratio":     fanRatio,
		"fan_runtime_hours":  s.FanRuntime(),
		"alarm_state":        s.AlarmState(),
		"ammonia_peak":       s.AmmoniaPeak(),
		"ammonia_average":    s.AmmoniaAverage(),
		"quota_window_sum":   s.QuotaWindowSum(),
		"quota_cycle_left":   s.QuotaCycleRemaining(),
		"valve_mode":         s.ValveMode(),
		"valve_target":       s.ValveTarget(),
		"defrost_interval":   s.DefrostInterval().String(),
		"defrost_duration":   s.DefrostDuration().String(),
		"registry_keys":      len(s.storeRegistry.Keys()),
	}
}

func (s *Service) SeriesCount() int {
	return s.series.Count()
}

func (s *Service) SeriesAverage() float64 {
	return s.series.Average()
}

func (s *Service) RoomTarget(roomID string) (float64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	controller, ok := s.roomStates[roomID]
	if !ok {
		return 0, false
	}
	return controller.Target(), true
}

func (s *Service) SetRoomTarget(roomID string, target float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	controller, ok := s.roomStates[roomID]
	if !ok {
		return false
	}
	controller.SetTarget(target)
	return true
}

func (s *Service) RoomDeviation(roomID string) (float64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	roomValue, ok := s.rooms[roomID]
	if !ok {
		return 0, false
	}
	controller, ok := s.roomStates[roomID]
	if !ok {
		return 0, false
	}
	return controller.Deviation(roomValue.TemperatureValue()), true
}

func (s *Service) RemoveOccupancy(roomID string, amount float64) (float64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.occupancy[roomID]
	if !ok {
		return 0, false
	}
	return value.Remove(amount), true
}

func (s *Service) UpdateCapacity(tons float64, setpoint float64) {
	s.capacity.Update(tons, setpoint)
}

func (s *Service) CoolingStart() {
	s.cooling.Start()
}

func (s *Service) CoolingStop() {
	s.cooling.Stop()
}

func (s *Service) CoolingAddHours(hours float64) {
	s.cooling.AddHours(hours)
}

func (s *Service) AddFanDuty(active float64, total float64) {
	s.fanDuty.Add(active, total)
}

func (s *Service) AddFanRuntime(hours float64) {
	s.fanRuntime.Add(hours)
}

func (s *Service) FanRuntimeDue(interval float64) bool {
	return s.fanRuntime.Due(interval)
}

func (s *Service) AddAmmoniaLevel(value float64) {
	s.ammoniaLevels.Add(value)
}

func (s *Service) AmmoniaLevelCount() int {
	return s.ammoniaLevels.Count()
}

func (s *Service) AmmoniaMargin() float64 {
	return s.ammoniaMonitor.Margin()
}

func (s *Service) RecordQuotaWindow(amount int) {
	s.quotaWindow.Record(amount)
}

func (s *Service) QuotaWindowAllowed(amount int) bool {
	return s.quotaWindow.Allowed(amount)
}

func (s *Service) QuotaWindowSize() int {
	return s.quotaWindow.Size()
}

func (s *Service) ConsumeQuotaCycle(amount int) bool {
	return s.quotaCycle.Consume(amount)
}

func (s *Service) ResetQuotaCycle() {
	s.quotaCycle.Reset()
}

func (s *Service) RegistryGet(key string) (any, bool) {
	return s.storeRegistry.Get(key)
}

func (s *Service) RegistryKeys() []string {
	return s.storeRegistry.Keys()
}

func (s *Service) ZoneRooms(zone string) []string {
	return s.membership.List(zone)
}

func (s *Service) RemoveRoomFromZone(roomID string) {
	s.membership.Remove(roomID)
}

func (s *Service) SetValveTarget(target float64) {
	s.actuator.SetTarget(target)
}

func (s *Service) ValveTarget() float64 {
	return s.actuator.Target()
}

func (s *Service) DefrostInterval() time.Duration {
	return s.defrostSchedule.Interval()
}

func (s *Service) DefrostDuration() time.Duration {
	return s.defrostSchedule.Duration()
}

func (s *Service) ScheduleHighestBand(load float64) int {
	return s.schedule.HighestBand(load)
}

func (s *Service) ScheduleBandCount() int {
	return s.schedule.BandCount()
}

func (s *Service) StepBand(amount int) int {
	return s.band.Step(amount)
}
