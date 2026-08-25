package temp

import "sync"

type Series struct {
	mu      sync.Mutex
	samples []Sample
}

func NewSeries(capacity int) *Series {
	return &Series{samples: make([]Sample, 0, capacity)}
}

func (s *Series) Add(sample Sample) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples = append(s.samples, sample)
}

func (s *Series) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.samples)
}

func (s *Series) Average() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return 0
	}
	var total float64
	for _, sample := range s.samples {
		total += sample.Smoothed
	}
	return total / float64(len(s.samples))
}

func (s *Series) Slope() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) < 2 {
		return 0
	}
	first := s.samples[0].Smoothed
	last := s.samples[len(s.samples)-1].Smoothed
	return (last - first) / float64(len(s.samples)-1)
}
