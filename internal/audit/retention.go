package audit

func (l *Log) Trim(limit int) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.records) <= limit {
		return 0
	}
	removed := len(l.records) - limit
	l.records = append([]Record(nil), l.records[removed:]...)
	return removed
}

func (l *Log) Recent(n int) []Record {
	l.mu.Lock()
	defer l.mu.Unlock()
	if n <= 0 || n > len(l.records) {
		n = len(l.records)
	}
	result := make([]Record, n)
	copy(result, l.records[len(l.records)-n:])
	return result
}
