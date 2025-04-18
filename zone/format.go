package zone

import "time"

func (z *zone) Format() string {
	return z.value.Format(time.RFC3339)
}
