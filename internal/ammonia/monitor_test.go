package ammonia

import (
	"sync"
	"testing"
)

// TestEvaluateCalibrateSequential confirms the core invariant after a
// calibration settles: Evaluate(reading) trips exactly when reading >= the
// threshold that was just published, including the tight-threshold case where a
// reading below the old (looser) threshold but above the new (tighter) one
// must now trip. This is the deterministic behavioral assertion for the bug
// where an over-limit concentration could miss the trip because the
// evaluation used a stale, pre-calibration threshold.
func TestEvaluateCalibrateSequential(t *testing.T) {
	m := NewMonitor(30)

	// reading 28 is below the initial threshold 30 -> no trip.
	if m.Evaluate(28) {
		t.Fatalf("Evaluate(28) with threshold 30: expected no trip")
	}
	// Tighten to 25; 28 now exceeds it and must trip immediately.
	m.Calibrate(25)
	if !m.Evaluate(28) {
		t.Fatalf("Evaluate(28) with threshold 25: expected trip")
	}
	// Loosen back to 30; 28 must again not trip.
	m.Calibrate(30)
	if m.Evaluate(28) {
		t.Fatalf("Evaluate(28) with threshold 30: expected no trip")
	}
}

// TestEvaluateCalibrateConcurrent stress-tests the publish (Calibrate) vs
// read/evaluate path under the race detector. Before the fix, Calibrate wrote
// m.threshold without a lock while Evaluate/Threshold read it, so the race
// detector flagged a data race and an evaluation could observe a stale
// threshold. With the fix all accesses to m.threshold go through m.mu, so the
// race detector is clean and the atomic snapshot stays self-consistent.
func TestEvaluateCalibrateConcurrent(t *testing.T) {
	m := NewMonitor(30)

	var wg sync.WaitGroup
	const iterations = 2000

	// Publisher: oscillate the threshold between a loose (30) and a tight
	// (25) value.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			if i%2 == 0 {
				m.Calibrate(30)
			} else {
				m.Calibrate(25)
			}
		}
	}()

	// Evaluator: publish reading 28 and judge. 28 is between the two
	// thresholds, so the correct trip result depends entirely on which
	// threshold is in effect; we only assert the evaluation never panics and
	// the race detector stays clean (the deterministic invariant is covered
	// by TestEvaluateCalibrateSequential).
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = m.Evaluate(28)
		}
	}()

	// Snapshot reader: concurrently read the reading/threshold pair, which
	// must always be a consistent pair captured under the lock.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_, _ = m.ReadingAndThreshold()
		}
	}()

	wg.Wait()

	// After everything settles, the final state must be self-consistent:
	// the atomic snapshot's threshold matches a standalone threshold read,
	// and Evaluate is consistent with that threshold.
	_, finalThreshold := m.ReadingAndThreshold()
	if got := m.Threshold(); got != finalThreshold {
		t.Fatalf("threshold snapshot mismatch: ReadingAndThreshold=%v Threshold=%v", finalThreshold, got)
	}
	reading := finalThreshold + 1
	if !m.Evaluate(reading) {
		t.Fatalf("Evaluate(%.0f) with threshold %.0f: expected trip", reading, finalThreshold)
	}
}
