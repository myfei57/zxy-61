package room

type Room struct {
	ID          string
	Name        string
	Temperature float64
}

func NewRoom(id string, name string, temperature float64) *Room {
	return &Room{ID: id, Name: name, Temperature: temperature}
}

func (r *Room) SetTemperature(temperature float64) {
	r.Temperature = temperature
}

func (r *Room) TemperatureValue() float64 {
	return r.Temperature
}
