package comp

import (
	"sync"
	"testing"

	"coldstore/internal/evap"
)

// TestClaimLeadAtomicConcurrency mirrors the field incident: defrost end and
// a load command land together, both units race through ClaimLead. Before the
// fix ClaimLead did an unlocked Leader()/SetLeader() pair so both banks
// believed they were leader and both went running. Exactly one bank must win.
func TestClaimLeadAtomicConcurrency(t *testing.T) {
	lead := evap.NewLead()
	a := NewBank("comp-a", 3)
	b := NewBank("comp-b", 3)

	var wg sync.WaitGroup
	results := make(chan bool, 2)
	start := make(chan struct{})
	wg.Add(2)
	run := func(bk *Bank) {
		defer wg.Done()
		<-start
		results <- bk.ClaimLead(lead)
	}
	go run(a)
	go run(b)
	close(start)
	wg.Wait()
	close(results)

	var wins int
	for r := range results {
		if r {
			wins++
		}
	}
	if wins != 1 {
		t.Fatalf("expected exactly one bank to win the lead, got %d", wins)
	}
	if !a.Running() && !b.Running() {
		t.Fatal("winning bank should be running")
	}
	if a.Running() && b.Running() {
		t.Fatal("both banks running — both believed they were leader")
	}
}
