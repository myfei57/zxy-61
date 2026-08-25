package room

type State string

const (
	StateNormal   State = "normal"
	StateCooling  State = "cooling"
	StateDefrost  State = "defrost"
)

type Controller struct {
	state  State
	target float64
}

func NewController(target float64) *Controller {
	return &Controller{state: StateNormal, target: target}
}

func (c *Controller) State() State {
	return c.state
}

func (c *Controller) SetState(state State) {
	c.state = state
}

func (c *Controller) Target() float64 {
	return c.target
}

func (c *Controller) SetTarget(target float64) {
	c.target = target
}

func (c *Controller) Deviation(current float64) float64 {
	return current - c.target
}

func (c *Controller) NeedsCooling(current float64, deadband float64) bool {
	return current-c.target > deadband
}
