package verifycase

import (
	"sync"
	"testing"

	"coldstore/internal/comp"
	"coldstore/internal/room"
)

func TestCsCoolingLoadBatchRace(t *testing.T) {
	load := room.NewLoad()
	load.SetBaseline("room-1", 1)
	bank := comp.NewBank("comp-a", 3)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				load.StockIn("room-1", 0.01)
			}
		}()
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				value := load.Baseline("room-1")
				_ = bank.StageByLoad(value)
			}
		}()
	}
	wg.Wait()
}
