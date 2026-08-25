package valve

func (v *SupplyValve) Open() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.open = true
	v.position = 1
	return nil
}
