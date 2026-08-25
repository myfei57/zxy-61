package evap

import "sync"

type Baseline struct {
	mu     sync.Mutex
	rating float64
	revision int64
}

func NewBaseline(rating float64) *Baseline {
	return &Baseline{rating: rating}
}

func (b *Baseline) UpdateRating(rating float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.rating = rating
	b.revision++
}

func (b *Baseline) Rating() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.rating
}

func (b *Baseline) Revision() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.revision
}
