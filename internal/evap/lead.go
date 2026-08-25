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
	return l.leader
}

func (l *Lead) Standby() string {
	return l.standby
}

func (l *Lead) SetLeader(candidate string) {
	l.leader = candidate
}
