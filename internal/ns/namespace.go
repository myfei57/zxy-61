package ns

type Namespace struct {
	ID      string
	Name    string
	RoomIDs []string
}

func NewNamespace(id string, name string) *Namespace {
	return &Namespace{ID: id, Name: name, RoomIDs: []string{}}
}

func (n *Namespace) AddRoom(roomID string) {
	for _, existing := range n.RoomIDs {
		if existing == roomID {
			return
		}
	}
	n.RoomIDs = append(n.RoomIDs, roomID)
}
