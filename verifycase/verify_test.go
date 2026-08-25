package verifycase

import (
	"testing"

	"coldstore/internal/evap"
	"coldstore/internal/fan"
)

func TestCsEvapFanBaselineRevisionFresh(t *testing.T) {
	baseline := evap.NewBaseline(40)
	verdict := fan.NewVerdict(baseline)
	baseline.UpdateRating(60)
	if verdict.Judge(55) {
		t.Fatal("fan verdict must follow the installed rating")
	}
}
