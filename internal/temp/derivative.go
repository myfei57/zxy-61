package temp

import "sync"

type Derivative struct {
	mu      sync.Mutex
	last    Sample
	ready   bool
	value   float64
}

func NewDerivative() *Derivative {
	return &Derivative{}
}

func (d *Derivative) Add(sample Sample) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ready {
		d.value = sample.Smoothed - d.last.Smoothed
	}
	d.last = sample
	d.ready = true
}

func (d *Derivative) Value() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.value
}
