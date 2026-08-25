package comp

import "sync"

type Schedule struct {
	mu    sync.Mutex
	bands []float64
}

func NewSchedule(bands []float64) *Schedule {
	copied := make([]float64, len(bands))
	copy(copied, bands)
	return &Schedule{bands: copied}
}

func (s *Schedule) Band(load float64) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	band := 1
	for _, threshold := range s.bands {
		if load >= threshold {
			band++
		}
	}
	return band
}

func (s *Schedule) BandCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.bands) + 1
}

func (s *Schedule) HighestBand(load float64) int {
	return s.Band(load)
}
