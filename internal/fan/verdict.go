package fan

import (
	"sync"

	"coldstore/internal/evap"
)

type Verdict struct {
	mu     sync.Mutex
	cached float64
}

func NewVerdict(baseline *evap.Baseline) *Verdict {
	return &Verdict{cached: baseline.Rating()}
}

func (v *Verdict) Judge(temperature float64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return temperature > v.cached
}
