package zone

import "time"

type zone struct {
	value time.Time
	tz    *time.Location
}

type Zone interface {
	Format() string
	Value() time.Time
}
