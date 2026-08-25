package quota

func (t *Tracker) Reserve(amount int) bool {
	return t.Allow(amount)
}
