package valve

func (v *SupplyValve) SetPosition(position float64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.position = position
}
