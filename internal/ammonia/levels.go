package ammonia

import "sync"

type Levels struct {
	mu      sync.Mutex
	history []float64
}

func NewLevels(capacity int) *Levels {
	return &Levels{history: make([]float64, 0, capacity)}
}

func (l *Levels) Add(value float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.history = append(l.history, value)
}

func (l *Levels) Peak() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	peak := 0.0
	for _, value := range l.history {
		if value > peak {
			peak = value
		}
	}
	return peak
}

func (l *Levels) Average() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.history) == 0 {
		return 0
	}
	var total float64
	for _, value := range l.history {
		total += value
	}
	return total / float64(len(l.history))
}

func (l *Levels) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.history)
}
