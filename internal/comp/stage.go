package comp

import (
	"errors"

	"coldstore/internal/evap"
)

func (b *Bank) StageLoad(pump *evap.OilPump) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if pump == nil {
		pump = b.pump
	}
	if b.loaded >= b.stages {
		return errors.New("all stages already loaded")
	}
	b.loaded++
	b.sequence = append(b.sequence, "stage-load")
	b.running = true
	if err := pump.Run(); err != nil {
		return err
	}
	b.sequence = append(b.sequence, "pump-run")
	return nil
}

func (b *Bank) StageByLoad(load float64) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	if load >= 4 {
		return b.ensureStage(3)
	}
	if load >= 2 {
		return b.ensureStage(2)
	}
	return b.ensureStage(1)
}
