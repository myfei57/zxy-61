package verifycase

import (
	"sync"
	"testing"

	"coldstore/internal/ammonia"
	"coldstore/internal/comp"
)

func TestCsAmmoniaThresholdPublishRace(t *testing.T) {
	monitor := ammonia.NewMonitor(30)
	interlock := comp.NewInterlock(monitor)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				monitor.Calibrate(float64(10 + (i+j)%20))
			}
		}(i)
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = interlock.Trip(12)
			}
		}()
	}
	wg.Wait()
}
