package fan

import "sync"

type Fan struct {
	mu      sync.Mutex
	running bool
}

func NewFan() *Fan {
	return &Fan{}
}

func (f *Fan) Stop() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.running = false
}

func (f *Fan) Running() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.running
}
