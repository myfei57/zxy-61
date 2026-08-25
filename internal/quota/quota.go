package quota

import "sync"

type Tracker struct {
	mu    sync.Mutex
	limit int
	used  int
}

func NewTracker(limit int) *Tracker {
	return &Tracker{limit: limit}
}

func (t *Tracker) Allow(amount int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.used+amount > t.limit {
		return false
	}
	t.used += amount
	return true
}

func (t *Tracker) Used() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.used
}

func (t *Tracker) Limit() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.limit
}

func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.used = 0
}
