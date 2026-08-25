package evap

import "sync"

type OilPump struct {
	mu      sync.Mutex
	running bool
}

func NewOilPump() *OilPump {
	return &OilPump{}
}

func (p *OilPump) Run() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.running = true
	return nil
}

func (p *OilPump) Running() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}
