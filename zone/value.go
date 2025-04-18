package zone

import "time"

func (z *zone) Value() time.Time {
	return z.value
}
