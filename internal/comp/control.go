package comp

import (
	"sync"

	"coldstore/internal/room"
	"coldstore/internal/temp"
)

func (b *Bank) Control(window *temp.Window, threshold float64) int {
	needMore := window.Verdict(threshold)
	b.mu.Lock()
	defer b.mu.Unlock()
	if needMore {
		return b.ensureStage(b.stages)
	}
	return b.ensureStage(1)
}

type RoomDriver struct {
	mu      sync.Mutex
	mapping *room.Mapping
	cached  map[string]string
}

func NewRoomDriver(mapping *room.Mapping) *RoomDriver {
	return &RoomDriver{mapping: mapping, cached: mapping.Snapshot()}
}

func (d *RoomDriver) RoomForMachine(machineID string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for roomID, machine := range d.cached {
		if machine == machineID {
			return roomID, true
		}
	}
	return "", false
}
