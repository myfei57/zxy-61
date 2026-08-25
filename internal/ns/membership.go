package ns

import "sync"

type Membership struct {
	mu       sync.RWMutex
	zones    map[string][]string
	roomZone map[string]string
}

func NewMembership() *Membership {
	return &Membership{zones: map[string][]string{}, roomZone: map[string]string{}}
}

func (m *Membership) Assign(zone string, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if previous, ok := m.roomZone[roomID]; ok && previous != zone {
		m.removeFromZone(previous, roomID)
	}
	m.roomZone[roomID] = zone
	m.zones[zone] = append(m.zones[zone], roomID)
}

func (m *Membership) ZoneOf(roomID string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	zone, ok := m.roomZone[roomID]
	return zone, ok
}

func (m *Membership) List(zone string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]string, len(m.zones[zone]))
	copy(result, m.zones[zone])
	return result
}

func (m *Membership) Remove(roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if zone, ok := m.roomZone[roomID]; ok {
		m.removeFromZone(zone, roomID)
		delete(m.roomZone, roomID)
	}
}

func (m *Membership) removeFromZone(zone string, roomID string) {
	filtered := m.zones[zone][:0]
	for _, id := range m.zones[zone] {
		if id != roomID {
			filtered = append(filtered, id)
		}
	}
	m.zones[zone] = filtered
}
