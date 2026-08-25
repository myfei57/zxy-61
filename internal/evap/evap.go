package evap

type Condenser struct {
	primary  string
	standby  string
	pump     *OilPump
	lead     *Lead
	baseline *Baseline
}

func NewCondenser(primary string, standby string, rating float64) *Condenser {
	lead := NewLead()
	lead.Seed(primary, standby)
	return &Condenser{
		primary:  primary,
		standby:  standby,
		pump:     NewOilPump(),
		lead:     lead,
		baseline: NewBaseline(rating),
	}
}

func (c *Condenser) Pump() *OilPump {
	return c.pump
}

func (c *Condenser) Lead() *Lead {
	return c.lead
}

func (c *Condenser) Baseline() *Baseline {
	return c.baseline
}
