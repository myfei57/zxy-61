package room

import "sync"

type Mapping struct {
	mu            sync.RWMutex
	roomToMachine map[string]string
	machineToRoom map[string]string
	revision      int64
}

func NewMapping() *Mapping {
	return &Mapping{
		roomToMachine: map[string]string{},
		machineToRoom: map[string]string{},
	}
}

func (m *Mapping) Set(roomID string, machineID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if previous, ok := m.machineToRoom[machineID]; ok {
		delete(m.roomToMachine, previous)
	}
	if previous, ok := m.roomToMachine[roomID]; ok {
		delete(m.machineToRoom, previous)
	}
	m.roomToMachine[roomID] = machineID
	m.machineToRoom[machineID] = roomID
	m.revision++
}

func (m *Mapping) RoomForMachine(machineID string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.machineToRoom[machineID]
	return value, ok
}

func (m *Mapping) Renumber(oldRoomID string, newRoomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	machineID, ok := m.roomToMachine[oldRoomID]
	if !ok {
		return
	}
	delete(m.roomToMachine, oldRoomID)
	m.roomToMachine[newRoomID] = machineID
	m.machineToRoom[machineID] = newRoomID
	m.revision++
}

func (m *Mapping) Revision() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.revision
}

func (m *Mapping) Snapshot() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]string, len(m.roomToMachine))
	for roomID, machineID := range m.roomToMachine {
		result[roomID] = machineID
	}
	return result
}
