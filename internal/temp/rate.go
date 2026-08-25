package temp

import "sync"

type Window struct {
	mu      sync.Mutex
	samples []Sample
}

func NewWindow(capacity int) *Window {
	return &Window{samples: make([]Sample, 0, capacity)}
}

func (w *Window) Add(sample Sample) {
	w.samples = append(w.samples, sample)
}

func (w *Window) delta(raw bool) float64 {
	if len(w.samples) < 2 {
		return 0
	}
	first := w.samples[0]
	last := w.samples[len(w.samples)-1]
	if raw {
		return last.Raw - first.Raw
	}
	return last.Smoothed - first.Smoothed
}

func (w *Window) Verdict(threshold float64) bool {
	return w.delta(false) >= threshold
}

func (w *Window) RawVerdict(threshold float64) bool {
	return w.delta(true) >= threshold
}
