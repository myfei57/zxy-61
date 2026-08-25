package service

import (
	"coldstore/internal/audit"
	"coldstore/internal/room"
)

func newRoom(id string, name string, temperature float64) *room.Room {
	return room.NewRoom(id, name, temperature)
}

func (s *Service) ListAudit() []audit.Record {
	return s.auditLog.Recent(100)
}

func (s *Service) ExportAudit() ([]byte, error) {
	return s.auditLog.Export()
}
