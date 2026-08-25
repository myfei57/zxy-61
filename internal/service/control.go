package service

import (
	"fmt"

	"coldstore/internal/room"
)

func (s *Service) AddRoom(id string, name string, temperature float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id == "" {
		return fmt.Errorf("room id is required")
	}
	if _, exists := s.rooms[id]; exists {
		return fmt.Errorf("room already exists: %s", id)
	}
	roomValue := newRoom(id, name, temperature)
	s.rooms[id] = roomValue
	s.mapping.Set(id, "comp-"+id)
	if namespace, ok := s.namespaces.Get("main"); ok {
		namespace.AddRoom(id)
	}
	s.roomStates[id] = room.NewController(temperature)
	s.occupancy[id] = room.NewOccupancy(0)
	s.membership.Assign("main", id)
	s.coolingLoad.SetBaseline(id, 1)
	s.auditLog.Add("console", "add-room", id)
	return nil
}

func (s *Service) SetRoomTemperature(id string, temperature float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	roomValue, ok := s.rooms[id]
	if !ok {
		return fmt.Errorf("room not found: %s", id)
	}
	roomValue.SetTemperature(temperature)
	sample := s.sampler.Read(temperature)
	s.window.Add(sample)
	s.series.Add(sample)
	s.derivative.Add(sample)
	if controller, ok := s.roomStates[id]; ok && controller.NeedsCooling(temperature, 1.0) {
		controller.SetState(room.StateCooling)
	}
	s.snapshot.Set("room:"+id, temperature)
	s.auditLog.Add("console", "set-temperature", id)
	return nil
}

func (s *Service) StockIn(id string, added float64) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rooms[id]; !ok {
		return 0, fmt.Errorf("room not found: %s", id)
	}
	baseline := s.coolingLoad.StockIn(id, added)
	if occupancy, ok := s.occupancy[id]; ok {
		occupancy.Add(added)
	}
	s.auditLog.Add("console", "stock-in", id)
	return baseline, nil
}

func (s *Service) StageByLoad(id string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rooms[id]; !ok {
		return 0, fmt.Errorf("room not found: %s", id)
	}
	load := s.coolingLoad.Baseline(id)
	stages := s.bank.StageByLoad(load)
	s.band.Set(s.schedule.Band(load))
	s.auditLog.Add("comp", "stage-by-load", id)
	return stages, nil
}

func (s *Service) Reconcile(threshold float64) int {
	stages := s.bank.Control(s.window, threshold)
	s.auditLog.Add("comp", "reconcile", "window")
	return stages
}

func (s *Service) RoomForMachine(machineID string) (string, bool) {
	return s.roomDriver.RoomForMachine(machineID)
}

func (s *Service) StartPlant() error {
	if err := s.supplyValve.Open(); err != nil {
		return err
	}
	if err := s.bank.Start(s.supplyValve); err != nil {
		return err
	}
	s.auditLog.Add("console", "start-plant", "start")
	return nil
}

func (s *Service) StageLoad() error {
	if err := s.bank.StageLoad(s.condenser.Pump()); err != nil {
		return err
	}
	s.auditLog.Add("console", "stage-load", "stage")
	return nil
}
