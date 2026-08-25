package comp

import (
	"sync"

	"coldstore/internal/evap"
)

type Bank struct {
	mu       sync.Mutex
	id       string
	stages   int
	loaded   int
	running  bool
	pump     *evap.OilPump
	sequence []string
}

func NewBank(id string, stages int) *Bank {
	return &Bank{
		id:       id,
		stages:   stages,
		pump:     evap.NewOilPump(),
		sequence: []string{},
	}
}

func (b *Bank) StageCount() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.stages
}

func (b *Bank) Loaded() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.loaded
}

func (b *Bank) Running() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.running
}

func (b *Bank) Pump() *evap.OilPump {
	return b.pump
}

func (b *Bank) SequenceEvents() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	result := make([]string, len(b.sequence))
	copy(result, b.sequence)
	return result
}

func (b *Bank) ensureStage(limit int) int {
	for b.loaded < limit && b.loaded < b.stages {
		b.loaded++
	}
	if b.loaded > 0 {
		b.running = true
	}
	return b.loaded
}
