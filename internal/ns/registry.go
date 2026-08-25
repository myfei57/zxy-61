package ns

import "sync"

type Registry struct {
	mu         sync.RWMutex
	namespaces map[string]*Namespace
}

func NewRegistry() *Registry {
	return &Registry{namespaces: map[string]*Namespace{}}
}

func (r *Registry) Put(n *Namespace) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.namespaces[n.ID] = n
}

func (r *Registry) Get(id string) (*Namespace, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.namespaces[id]
	return value, ok
}
