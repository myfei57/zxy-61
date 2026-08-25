package defrost

func (c *Controller) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fan != nil {
		if err := c.fan.StopAndPersist(); err != nil {
			return err
		}
		c.steps = append(c.steps, "fan-stop-persisted")
	}
	c.heater = true
	c.steps = append(c.steps, "heater-start")
	return nil
}
