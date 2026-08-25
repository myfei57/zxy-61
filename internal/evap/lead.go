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

// Claim atomically takes the leader role for candidate when the role is
// vacant or already held by candidate. The vacancy check and the write
// happen under a single lock hold so that two units racing to become
// leader at the same instant (e.g. defrost end colliding with a load
// command) cannot both observe an empty role and both win. Returns true
// when candidate is (or has become) the leader, false when another unit
// already holds it.
func (l *Lead) Claim(candidate string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.leader != "" && l.leader != candidate {
		return false
	}
	l.leader = candidate
	return true
}

// Acquire is the first-wins election: candidate becomes leader only when
// the role is empty; otherwise it is recorded as standby and rejected.
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

// SetLeader forces candidate into the leader role. It is intended for
// administrative fail-over; the election path must use Claim so that the
// vacancy check and the assignment stay atomic.
func (l *Lead) SetLeader(candidate string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.leader = candidate
}
