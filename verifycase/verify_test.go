package verifycase

import (
	"sync"
	"testing"

	"coldstore/internal/comp"
	"coldstore/internal/evap"
)

func TestCsConcurrentCompArbitration(t *testing.T) {
	lead := evap.NewLead()
	bankA := comp.NewBank("comp-a", 1)
	bankB := comp.NewBank("comp-b", 1)
	var results [2]bool
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		results[0] = bankA.ClaimLead(lead)
	}()
	go func() {
		defer wg.Done()
		results[1] = bankB.ClaimLead(lead)
	}()
	wg.Wait()
	leaders := 0
	for _, result := range results {
		if result {
			leaders++
		}
	}
	if leaders != 1 {
		t.Fatalf("exactly one compressor must become lead, got %d", leaders)
	}
}
