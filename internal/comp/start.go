package comp

import (
	"coldstore/internal/evap"
	"coldstore/internal/valve"
)

func (b *Bank) Start(v *valve.SupplyValve) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if v != nil {
		if err := v.Open(); err != nil {
			return err
		}
	}
	b.running = true
	return nil
}

func (b *Bank) ClaimLead(lead *evap.Lead) bool {
	if !lead.Acquire(b.id) {
		return false
	}
	b.mu.Lock()
	b.running = true
	b.mu.Unlock()
	return true
}
