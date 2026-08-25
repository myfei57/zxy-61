package valve

import "sync"

type Mode string

const (
	ModeAuto    Mode = "auto"
	ModeManual  Mode = "manual"
	ModeDefrost Mode = "defrost"
)

type Actuator struct {
	mu     sync.Mutex
	mode   Mode
	target float64
}

func NewActuator(mode Mode, target float64) *Actuator {
	return &Actuator{mode: mode, target: target}
}

func (a *Actuator) Mode() Mode {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.mode
}

func (a *Actuator) SetTarget(target float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.target = target
}

func (a *Actuator) Target() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.target
}
