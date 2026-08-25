package comp

import "sync"

type Band struct {
	mu      sync.Mutex
	current int
}

func NewBand(initial int) *Band {
	return &Band{current: initial}
}

func (b *Band) Current() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.current
}

func (b *Band) Set(value int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = value
}

func (b *Band) Step(amount int) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current += amount
	if b.current < 1 {
		b.current = 1
	}
	return b.current
}
