package defrost

import "coldstore/internal/valve"

type RunController struct {
	running bool
}

func NewRunController() *RunController {
	return &RunController{}
}

func (r *RunController) Run(v *valve.SupplyValve) error {
	r.running = true
	if v != nil {
		v.SetPosition(0)
	}
	return nil
}
