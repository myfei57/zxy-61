package service

import (
	"testing"
)

// TestCheckInterlockCalibrateSequential confirms the behavioral fix at the
// service level: once a tighter threshold is published by CalibrateAmmonia, a
// reading that was previously below the old (looser) threshold but now exceeds
// the new (tighter) one must trip immediately — not against a stale,
// pre-calibration threshold. This is the deterministic regression test for the
// reported bug where an over-limit concentration missed the interlock because
// the evaluation used the pre-calibration threshold.
func TestCheckInterlockCalibrateSequential(t *testing.T) {
	srv, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Initial threshold (from NewMonitor) is 30. A reading of 28 is below it
	// and must not trip.
	if srv.CheckInterlock(28) {
		t.Fatalf("CheckInterlock(28) with threshold 30: expected no trip")
	}
	// Tighten the threshold to 25. 28 now exceeds it and must trip.
	srv.CalibrateAmmonia(25)
	if !srv.CheckInterlock(28) {
		t.Fatalf("CheckInterlock(28) with threshold 25: expected trip")
	}
	// Loosen back to 30; 28 must again not trip.
	srv.CalibrateAmmonia(30)
	if srv.CheckInterlock(28) {
		t.Fatalf("CheckInterlock(28) with threshold 30: expected no trip")
	}
}

// TestCheckInterlockPublishesReading confirms CheckInterlock records the
// supplied reading so that the status snapshot reports a reading consistent
// with what was judged (i.e. the reading publication and the trip evaluation
// share a single atomic step).
func TestCheckInterlockPublishesReading(t *testing.T) {
	srv, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	srv.CheckInterlock(42)
	reading, threshold := srv.ammoniaMonitor.ReadingAndThreshold()
	if reading != 42 {
		t.Fatalf("reading not published: got %v want 42", reading)
	}
	if threshold != 30 {
		t.Fatalf("threshold changed unexpectedly: got %v want 30", threshold)
	}
}
