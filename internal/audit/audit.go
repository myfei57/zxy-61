package audit

import (
	"sync"

	"github.com/google/uuid"
)

type Log struct {
	mu      sync.Mutex
	records []Record
}

func NewLog() *Log {
	return &Log{records: []Record{}}
}

func (l *Log) Add(actor string, action string, detail string) Record {
	record := NewRecord(uuid.NewString(), actor, action, detail)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, record)
	return record
}

func (l *Log) List() []Record {
	l.mu.Lock()
	defer l.mu.Unlock()
	result := make([]Record, len(l.records))
	copy(result, l.records)
	return result
}

func (l *Log) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.records)
}
