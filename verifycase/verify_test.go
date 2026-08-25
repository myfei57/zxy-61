package verifycase

import (
	"testing"

	"coldstore/internal/comp"
	"coldstore/internal/room"
)

func TestCsRoomMappingRevisionFresh(t *testing.T) {
	mapping := room.NewMapping()
	mapping.Set("room-a", "machine-m")
	driver := comp.NewRoomDriver(mapping)
	mapping.Renumber("room-a", "room-b")
	got, ok := driver.RoomForMachine("machine-m")
	if !ok {
		t.Fatal("machine should be mapped")
	}
	if got != "room-b" {
		t.Fatalf("machine must follow the renumbered room, got %s", got)
	}
}
