package spotinstrumentclient

import "time"

type Market struct {
	id        string
	active    bool
	removedAt *time.Time
}

func NewMarket(id string, active bool, removedAt *time.Time) Market {
	return Market{id: id, active: active, removedAt: removedAt}
}

func (s Market) ID() string {
	return s.id
}
func (s Market) Active() bool {
	return s.active
}
func (s Market) RemovedAt() *time.Time {
	return s.removedAt
}
