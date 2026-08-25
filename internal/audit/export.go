package audit

import "encoding/json"

func (l *Log) Export() ([]byte, error) {
	records := l.List()
	return json.Marshal(records)
}
