package verifycase

import (
	"sync"
	"testing"

	"coldstore/internal/comp"
	"coldstore/internal/temp"
)

func TestCsTempWindowVerdictRace(t *testing.T) {
	window := temp.NewWindow(64)
	bank := comp.NewBank("comp-a", 3)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				window.Add(temp.Sample{Raw: float64(i + j), Smoothed: float64(i + j)})
			}
		}(i)
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = bank.Control(window, 10)
			}
		}()
	}
	wg.Wait()
}
