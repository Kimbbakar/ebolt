package eblot

import "time"

type cachePayload struct {
	Value     interface{}
	CreatedAt time.Time
	Exp       *time.Time
}

func (p cachePayload) isExpired() bool {
	if p.Exp != nil {
		now := time.Now()
		return now.After(*p.Exp)
	}

	return false
}
