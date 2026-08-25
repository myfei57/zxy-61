package verifycase

import (
	"testing"

	"coldstore/internal/comp"
	"coldstore/internal/evap"
)

func TestCsCompStageOrder(t *testing.T) {
	bank := comp.NewBank("comp-a", 1)
	pump := evap.NewOilPump()
	if err := bank.StageLoad(pump); err != nil {
		t.Fatalf("stage load failed: %v", err)
	}
	events := bank.SequenceEvents()
	if len(events) == 0 || events[0] != "pump-run" {
		t.Fatalf("oil pump must run before stage load, got %v", events)
	}
	if !pump.Running() {
		t.Fatal("oil pump should be running")
	}
}
