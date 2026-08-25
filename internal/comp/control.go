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
	mu       sync.Mutex
	mapping  *room.Mapping
	cached   map[string]string
	revision int64
}

func NewRoomDriver(mapping *room.Mapping) *RoomDriver {
	return &RoomDriver{mapping: mapping, cached: mapping.Snapshot(), revision: mapping.Revision()}
}

func (d *RoomDriver) RoomForMachine(machineID string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if revision := d.mapping.Revision(); revision != d.revision {
		d.cached = d.mapping.Snapshot()
		d.revision = revision
	}
	for roomID, machine := range d.cached {
		if machine == machineID {
			return roomID, true
		}
	}
	return "", false
}
