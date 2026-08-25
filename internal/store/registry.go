package store

import "sync"

type Registry struct {
	mu      sync.RWMutex
	entries map[string]any
}

func NewRegistry() *Registry {
	return &Registry{entries: map[string]any{}}
}

func (r *Registry) Put(key string, value any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[key] = value
}

func (r *Registry) Get(key string) (any, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.entries[key]
	return value, ok
}

func (r *Registry) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]string, 0, len(r.entries))
	for key := range r.entries {
		keys = append(keys, key)
	}
	return keys
}
