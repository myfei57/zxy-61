package quota

import "sync"

type Window struct {
	mu     sync.Mutex
	limit  int
	values []int
	cursor int
	filled int
}

func NewWindow(limit int, size int) *Window {
	return &Window{limit: limit, values: make([]int, size)}
}

func (w *Window) Record(amount int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.values[w.cursor] = amount
	w.cursor = (w.cursor + 1) % len(w.values)
	if w.filled < len(w.values) {
		w.filled++
	}
}

func (w *Window) Sum() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	total := 0
	for i := 0; i < w.filled; i++ {
		total += w.values[i]
	}
	return total
}

func (w *Window) Allowed(amount int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Sum()+amount <= w.limit
}

func (w *Window) Size() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.values)
}
