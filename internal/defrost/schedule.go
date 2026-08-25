package defrost

import (
	"sync"
	"time"
)

type Schedule struct {
	mu       sync.Mutex
	interval time.Duration
	duration time.Duration
	last     time.Time
}

func NewSchedule(interval time.Duration, duration time.Duration) *Schedule {
	return &Schedule{interval: interval, duration: duration}
}

func (s *Schedule) Interval() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.interval
}

func (s *Schedule) Duration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.duration
}
