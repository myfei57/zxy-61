package evap

import "sync"

type Lead struct {
	mu      sync.Mutex
	leader  string
	standby string
}

func NewLead() *Lead {
	return &Lead{}
}

func (l *Lead) Seed(primary string, standby string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.leader = primary
	l.standby = standby
}

func (l *Lead) Acquire(candidate string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.leader == "" {
		l.leader = candidate
		return true
	}
	if l.standby == "" && candidate != l.leader {
		l.standby = candidate
	}
	return false
}

func (l *Lead) Leader() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.leader
}

func (l *Lead) Standby() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.standby
}
