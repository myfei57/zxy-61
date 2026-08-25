package fan

import "sync"

type Duty struct {
	mu     sync.Mutex
	active float64
	total  float64
}

func NewDuty() *Duty {
	return &Duty{}
}

func (d *Duty) Add(active float64, total float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.active += active
	d.total += total
}

func (d *Duty) Ratio() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.total == 0 {
		return 0
	}
	return d.active / d.total
}

func (d *Duty) Active() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.active
}

func (d *Duty) Total() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.total
}
