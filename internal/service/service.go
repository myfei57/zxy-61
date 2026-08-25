package service

import (
	"sync"
	"time"

	"coldstore/internal/ammonia"
	"coldstore/internal/audit"
	"coldstore/internal/comp"
	"coldstore/internal/defrost"
	"coldstore/internal/evap"
	"coldstore/internal/fan"
	"coldstore/internal/ns"
	"coldstore/internal/quota"
	"coldstore/internal/room"
	"coldstore/internal/store"
	"coldstore/internal/temp"
	"coldstore/internal/valve"
)

type Service struct {
	mu                 sync.Mutex
	path               string
	snapshot           *store.Snapshot
	fileStore          *store.Store
	auditLog           *audit.Log
	namespaces         *ns.Registry
	rooms              map[string]*room.Room
	mapping            *room.Mapping
	coolingLoad        *room.Load
	sampler            *temp.Sampler
	window             *temp.Window
	bank               *comp.Bank
	roomDriver         *comp.RoomDriver
	interlock          *comp.Interlock
	condenser          *evap.Condenser
	fanUnit            *fan.Fan
	fanStop            *fan.StopUnit
	fanVerdict         *fan.Verdict
	defrostController  *defrost.Controller
	defrostRun         *defrost.RunController
	supplyValve        *valve.SupplyValve
	ammoniaMonitor     *ammonia.Monitor
	quotaTracker       *quota.Tracker
	roomStates         map[string]*room.Controller
	occupancy          map[string]*room.Occupancy
	series             *temp.Series
	derivative         *temp.Derivative
	schedule           *comp.Schedule
	band               *comp.Band
	capacity           *evap.Capacity
	cooling            *evap.Cooling
	fanDuty            *fan.Duty
	fanRuntime         *fan.Runtime
	ammoniaLevels      *ammonia.Levels
	quotaWindow        *quota.Window
	quotaCycle         *quota.Cycle
	membership         *ns.Membership
	actuator           *valve.Actuator
	defrostSchedule    *defrost.Schedule
	storeRegistry      *store.Registry
}

func New(path string) (*Service, error) {
	monitor := ammonia.NewMonitor(30)
	condenser := evap.NewCondenser("comp-a", "comp-b", 40)
	baseline := condenser.Baseline()
	bank := comp.NewBank("comp-a", 3)
	fanStop := fan.NewStopUnit()
	defaultNamespace := ns.NewNamespace("main", "main-rooms")
	registry := ns.NewRegistry()
	registry.Put(defaultNamespace)
	mapping := room.NewMapping()

	return &Service{
		path:               path,
		snapshot:           store.NewSnapshot(),
		fileStore:          store.New(path),
		auditLog:           audit.NewLog(),
		namespaces:         registry,
		rooms:              map[string]*room.Room{},
		mapping:            mapping,
		coolingLoad:        room.NewLoad(),
		sampler:            temp.NewSampler(0.7),
		window:             temp.NewWindow(128),
		bank:               bank,
		roomDriver:         comp.NewRoomDriver(mapping),
		interlock:          comp.NewInterlock(monitor),
		condenser:          condenser,
		fanUnit:            fan.NewFan(),
		fanStop:            fanStop,
		fanVerdict:         fan.NewVerdict(baseline),
		defrostController:  defrost.NewController(fanStop),
		defrostRun:         defrost.NewRunController(),
		supplyValve:        valve.NewSupplyValve(),
		ammoniaMonitor:     monitor,
		quotaTracker:       quota.NewTracker(1000),
		roomStates:         map[string]*room.Controller{},
		occupancy:          map[string]*room.Occupancy{},
		series:             temp.NewSeries(256),
		derivative:         temp.NewDerivative(),
		schedule:           comp.NewSchedule([]float64{2, 4}),
		band:               comp.NewBand(1),
		capacity:           evap.NewCapacity(20, -18),
		cooling:            evap.NewCooling(),
		fanDuty:            fan.NewDuty(),
		fanRuntime:         fan.NewRuntime(),
		ammoniaLevels:      ammonia.NewLevels(256),
		quotaWindow:        quota.NewWindow(1000, 10),
		quotaCycle:         quota.NewCycle(1000),
		membership:         ns.NewMembership(),
		actuator:           valve.NewActuator(valve.ModeAuto, 1),
		defrostSchedule:    defrost.NewSchedule(6*time.Hour, 20*time.Minute),
		storeRegistry:      store.NewRegistry(),
	}, nil
}

func (s *Service) Load() error {
	var payload map[string]any
	if err := s.fileStore.Load(&payload); err != nil {
		return err
	}
	if payload != nil {
		s.snapshot.Set("loaded", true)
	}
	return nil
}

func (s *Service) Save() error {
	s.snapshot.Set("room_count", len(s.rooms))
	s.snapshot.Set("audit_count", s.auditLog.Count())
	s.auditLog.Trim(1000)
	return s.fileStore.SaveAtomic(s.snapshot.All())
}

func (s *Service) Health() map[string]any {
	return map[string]any{"status": "ok", "stage_count": s.bank.StageCount()}
}

func (s *Service) Status() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Snapshot the reading and threshold together so the reported pair is
	// consistent and cannot reflect a reading from one instant and a
	// threshold from another across a concurrent calibration.
	ammoniaReading, ammoniaThreshold := s.ammoniaMonitor.ReadingAndThreshold()
	return map[string]any{
		"room_count":         len(s.rooms),
		"compressor_running": s.bank.Running(),
		"compressor_loaded":  s.bank.Loaded(),
		"sequence":           s.bank.SequenceEvents(),
		"lead":               s.condenser.Lead().Leader(),
		"standby":            s.condenser.Lead().Standby(),
		"fan_running":        s.fanUnit.Running(),
		"oil_pump_running":   s.condenser.Pump().Running(),
		"fan_persisted":      s.fanStop.Persisted(),
		"heater_on":          s.defrostController.HeaterOn(),
		"defrost_steps":      s.defrostController.Steps(),
		"valve_open":         s.supplyValve.IsOpen(),
		"valve_position":     s.supplyValve.Position(),
		"ammonia_reading":    ammoniaReading,
		"ammonia_threshold":  ammoniaThreshold,
		"quota_used":         s.quotaTracker.Used(),
		"quota_limit":        s.quotaTracker.Limit(),
		"audit_count":        s.auditLog.Count(),
	}
}
