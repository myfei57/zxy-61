package fan

import (
	"sync"

	"coldstore/internal/evap"
)

type Verdict struct {
	mu       sync.Mutex
	baseline *evap.Baseline
	cached   float64
	revision int64
}

func NewVerdict(baseline *evap.Baseline) *Verdict {
	return &Verdict{baseline: baseline, cached: baseline.Rating(), revision: baseline.Revision()}
}

func (v *Verdict) Judge(temperature float64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	if revision := v.baseline.Revision(); revision != v.revision {
		v.cached = v.baseline.Rating()
		v.revision = revision
	}
	return temperature > v.cached
}
