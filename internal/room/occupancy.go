package room

import "sync"

type Occupancy struct {
	mu    sync.Mutex
	stock float64
}

func NewOccupancy(stock float64) *Occupancy {
	return &Occupancy{stock: stock}
}

func (o *Occupancy) Add(amount float64) float64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.stock += amount
	return o.stock
}

func (o *Occupancy) Remove(amount float64) float64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	if amount > o.stock {
		o.stock = 0
		return o.stock
	}
	o.stock -= amount
	return o.stock
}
