package verifycase

import (
	"sync"
	"testing"

	"coldstore/internal/defrost"
	"coldstore/internal/valve"
)

func TestCsConcurrentDefrostValve(t *testing.T) {
	supplyValve := valve.NewSupplyValve()
	runController := defrost.NewRunController()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_ = runController.Run(supplyValve)
	}()
	go func() {
		defer wg.Done()
		_ = supplyValve.Open()
	}()
	wg.Wait()
	position := supplyValve.Position()
	if position != 0 && position != 1 {
		t.Fatalf("valve position must be fully open or fully closed, got %v", position)
	}
}
