package fan

import "sync"

type Runtime struct {
	mu    sync.Mutex
	hours float64
}

func NewRuntime() *Runtime {
	return &Runtime{}
}

func (r *Runtime) Add(hours float64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hours += hours
}

func (r *Runtime) Hours() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.hours
}

func (r *Runtime) Due(interval float64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.hours >= interval
}
