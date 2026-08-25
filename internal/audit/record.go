package audit

import "time"

type Record struct {
	ID     string
	At     time.Time
	Actor  string
	Action string
	Detail string
}

func NewRecord(id string, actor string, action string, detail string) Record {
	return Record{ID: id, At: time.Now(), Actor: actor, Action: action, Detail: detail}
}
