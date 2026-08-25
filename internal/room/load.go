package room

import "sync"

type Load struct {
	mu       sync.RWMutex
	baseline map[string]float64
}

func NewLoad() *Load {
	return &Load{baseline: map[string]float64{}}
}

func (l *Load) Baseline(roomID string) float64 {
	return l.baseline[roomID]
}

func (l *Load) SetBaseline(roomID string, value float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.baseline[roomID] = value
}

func (l *Load) StockIn(roomID string, added float64) float64 {
	l.baseline[roomID] += added
	return l.baseline[roomID]
}
