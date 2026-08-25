package valve

import "sync"

type SupplyValve struct {
	mu       sync.Mutex
	open     bool
	position float64
}

func NewSupplyValve() *SupplyValve {
	return &SupplyValve{}
}

func (v *SupplyValve) IsOpen() bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.open
}

func (v *SupplyValve) Position() float64 {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.position
}
