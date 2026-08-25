package store

import "sync"

type Snapshot struct {
	mu     sync.RWMutex
	values map[string]any
}

func NewSnapshot() *Snapshot {
	return &Snapshot{values: map[string]any{}}
}

func (s *Snapshot) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
}

func (s *Snapshot) All() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]any, len(s.values))
	for key, value := range s.values {
		result[key] = value
	}
	return result
}
