package fan

import "sync"

type StopUnit struct {
	mu        sync.Mutex
	persisted bool
}

func NewStopUnit() *StopUnit {
	return &StopUnit{}
}

func (u *StopUnit) StopAndPersist() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.persisted = true
	return nil
}

func (u *StopUnit) Persisted() bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.persisted
}
