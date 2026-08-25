package ammonia

type AlarmState string

const (
	AlarmNormal  AlarmState = "normal"
	AlarmWarning AlarmState = "warning"
	AlarmTrip    AlarmState = "trip"
)

func (m *Monitor) State() AlarmState {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.reading >= m.threshold {
		return AlarmTrip
	}
	if m.reading >= m.threshold*0.8 {
		return AlarmWarning
	}
	return AlarmNormal
}

func (m *Monitor) Margin() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.threshold - m.reading
}
