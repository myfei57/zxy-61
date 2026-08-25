package temp

import "sync"

type Sample struct {
	Raw      float64
	Smoothed float64
}

type Sampler struct {
	mu    sync.Mutex
	last  float64
	alpha float64
}

func NewSampler(alpha float64) *Sampler {
	return &Sampler{alpha: alpha}
}

func (s *Sampler) Read(raw float64) Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	smoothed := raw
	if s.last != 0 {
		smoothed = s.alpha*raw + (1-s.alpha)*s.last
	}
	s.last = smoothed
	return Sample{Raw: raw, Smoothed: smoothed}
}
