package evap

import (
	"sync"
	"testing"
)

// TestClaimAtomicConcurrency reproduces the defrost-end / load-command race:
// two units call Claim at the same instant with an empty leader role. Before
// the fix the read-then-write in ClaimLead let both observe an empty role and
// both win, so two compressors loaded. Exactly one unit must become leader.
func TestClaimAtomicConcurrency(t *testing.T) {
	lead := NewLead()

	const units = 2
	var wg sync.WaitGroup
	winners := make(chan string, units)
	start := make(chan struct{})

	wg.Add(units)
	for i := 0; i < units; i++ {
		id := "comp-a"
		if i == 1 {
			id = "comp-b"
		}
		go func(id string) {
			defer wg.Done()
			<-start // fire both goroutines at the same instant
			if lead.Claim(id) {
				winners <- id
			}
		}(id)
	}
	close(start)
	wg.Wait()
	close(winners)

	var got []string
	for w := range winners {
		got = append(got, w)
	}
	if len(got) != 1 {
		t.Fatalf("expected exactly one leader, got %d (%v); election is not atomic", len(got), got)
	}
	if lead.Leader() == "" {
		t.Fatal("leader role is empty after election")
	}
}

// TestClaimIdempotentReentry: the leader may re-claim its own role (e.g. a
// load command arriving after it already won) without being rejected.
func TestClaimIdempotentReentry(t *testing.T) {
	lead := NewLead()
	if !lead.Claim("comp-a") {
		t.Fatal("first claim should win")
	}
	if !lead.Claim("comp-a") {
		t.Fatal("re-claim by the existing leader should succeed")
	}
	if lead.Leader() != "comp-a" {
		t.Fatalf("leader = %q, want comp-a", lead.Leader())
	}
}

// TestClaimRejectsChallenger: a different unit is refused once a leader exists.
func TestClaimRejectsChallenger(t *testing.T) {
	lead := NewLead()
	lead.Claim("comp-a")
	if lead.Claim("comp-b") {
		t.Fatal("a different unit must not be able to take an occupied leader role")
	}
	if lead.Leader() != "comp-a" {
		t.Fatalf("leader = %q, want comp-a", lead.Leader())
	}
}
